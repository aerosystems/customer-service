package app

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	rootCmd.Flags().BoolVar(&flags.Migration, "migration", false, "Run migration")
	rootCmd.Flags().BoolVar(&flags.RunHTTP, "http", false, "Run HTTP server")
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
	if flags.Migration {
		app := InitMigrationApp()
		if err := app.migration.Run(); err != nil {
			app.log.Fatalf("migration failed: %v", err)
		}
		app.log.Info("migration completed")
		return
	}

	app := InitServerApp()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	if flags.RunHTTP {
		group.Go(func() error {
			return app.http.Run()
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
