package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func (app *Server) handleSignals(ctx context.Context, cancel context.CancelFunc) error {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signalCh:
		// Handle graceful shutdown
		return app.gracefulShutdown(ctx, cancel)
	case <-ctx.Done():
		// Context cancelled, shutdown initiated elsewhere
		return nil
	}
}

func (app *Server) gracefulShutdown(ctx context.Context, cancel context.CancelFunc) error {
	cancel()
	app.log.Fatalf("app is shutting down")
	return app.http.Shutdown(ctx)
}
