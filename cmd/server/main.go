package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ttyobiwan/dstrat/internal/sqlite"
)

func getDB(getenv func(string) string) (*sql.DB, error) {
	dbName := getenv("DB_NAME")
	if dbName == "" {
		dbName = "default.sqlite"
	}

	db, err := sqlite.GetDB(dbName)
	if err != nil {
		return nil, fmt.Errorf("getting db: %v", err)
	}

	err = sqlite.Configure(db)
	if err != nil {
		return nil, fmt.Errorf("configuring db: %v", err)
	}

	err = sqlite.Migrate(db)
	if err != nil {
		return nil, fmt.Errorf("migrating db: %v", err)
	}

	return db, nil
}

func newServer(db *sql.DB) *echo.Echo {
	e := echo.New()

	e.Use(newLoggingMiddleware())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.RequestID())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{Timeout: time.Second * 25}))

	addRoutes(e, db)

	return e
}

func run(ctx context.Context, getenv func(string) string) error {
	// Set logger
	if getenv("DEBUG") != "true" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Get dependencies
	db, err := getDB(getenv)
	if err != nil {
		return err
	}

	// Start server
	srv := newServer(db)
	go func() {
		port := getenv("PORT")
		if port == "" {
			port = "8088"
		}
		if err := srv.Start(":" + port); err != nil {
			slog.Error("Error listening", "error", err)
		}
	}()

	// Prepare a proper shutdown
	done := make(chan struct{})
	go func() {
		<-ctx.Done()

		done <- struct{}{}
		slog.Info("Shutting down")
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error shutting down the server", "error", err)
		}
	}()
	<-done

	return nil
}

func main() {
	if err := run(context.Background(), os.Getenv); err != nil {
		slog.Error("Error starting the server", "error", err)
		os.Exit(1)
	}
}
