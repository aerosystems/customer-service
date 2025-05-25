package app

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/aerosystems/customer-service/internal/adapters"
	HTTPServer "github.com/aerosystems/customer-service/internal/ports/http"
)

func init() {
	rootCmd.Flags().BoolVar(&flags.Migration, "migration", false, "Run migration")
	rootCmd.Flags().BoolVar(&flags.RunHTTP, "http", false, "Run HTTP server")
}

type App struct {
	log        *logrus.Logger
	cfg        *Config
	migration  *adapters.Migration
	httpServer *HTTPServer.Server
}

func NewApp(
	log *logrus.Logger,
	cfg *Config,
	migration *adapters.Migration,
	httpServer *HTTPServer.Server,
) *App {
	return &App{
		log:        log,
		cfg:        cfg,
		migration:  migration,
		httpServer: httpServer,
	}
}

var rootCmd = &cobra.Command{
	Use:   "customer-service",
	Short: "Run customer-service",
	Run:   runRoot,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	app := InitApp()

	if flags.Migration {
		if err := app.migration.Run(); err != nil {
			app.log.Fatalf("migration failed: %v", err)
		}
		app.log.Info("migration completed")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	if flags.RunHTTP {
		group.Go(func() error {
			return app.httpServer.Run()
		})
	} else {
		cancel()
	}

	group.Go(func() error {
		return app.handleSignals(ctx, cancel)
	})

	if err := group.Wait(); err != nil {
		app.log.Errorf("error occurred: %v", err)
	}
}
