package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ttyobiwan/dstrat/internal/sqlite"
	"github.com/ttyobiwan/dstrat/internal/temporal"
	"github.com/ttyobiwan/dstrat/posts"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// Indicator that worker will be used for test purposes
const testQueueSentinel = "test"

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

func registerWorkflows(w worker.Worker, queue string, db *sql.DB) error {
	switch queue {
	case string(temporal.TemporalQueuePosts):
		workflows := posts.NewWorkflowManager(db)
		activities := posts.NewActivityManager(db)
		w.RegisterWorkflow(workflows.SendPostToTopicFollowers)
		w.RegisterWorkflow(workflows.SendPostBulk)
		w.RegisterWorkflow(workflows.SendPostMass)
		w.RegisterActivity(activities.GetPost)
		w.RegisterActivity(activities.GetFollowers)
		w.RegisterActivity(activities.SendSinglePost)
		w.RegisterActivity(activities.SendPostSequentially)
		w.RegisterActivity(activities.SendPostASequentially)
	// Worker for test purposes - register all workflows
	case testQueueSentinel:
		workflows := posts.NewWorkflowManager(db)
		activities := posts.NewActivityManager(db)
		w.RegisterWorkflow(workflows.SendPostToTopicFollowers)
		w.RegisterWorkflow(workflows.SendPostBulk)
		w.RegisterWorkflow(workflows.SendPostMass)
		w.RegisterActivity(activities.GetPost)
		w.RegisterActivity(activities.GetFollowers)
		w.RegisterActivity(activities.SendSinglePost)
		w.RegisterActivity(activities.SendPostSequentially)
		w.RegisterActivity(activities.SendPostASequentially)
	default:
		return fmt.Errorf("invalid queue name: %s", queue)
	}
	slog.Info("Registered workflows for the queue", "queue", queue)
	return nil
}

func run(getenv func(string) string, queue string) error {
	if getenv("DEBUG") != "true" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}

	hostPort := getenv("TEMPORAL_HOSTPORT")
	if hostPort == "" {
		hostPort = "127.0.0.1:7233"
	}
	c, err := client.Dial(client.Options{Logger: slog.Default(), HostPort: hostPort})
	if err != nil {
		return fmt.Errorf("creating client: %v", err)
	}
	defer c.Close()

	db, err := getDB(getenv)
	if err != nil {
		return err
	}
	defer db.Close()

	w := worker.New(c, queue, worker.Options{})
	if err := registerWorkflows(w, queue, db); err != nil {
		return fmt.Errorf("registering workflows: %v", err)
	}

	err = w.Run(worker.InterruptCh())
	if err != nil {
		return fmt.Errorf("starting worker: %v", err)
	}

	return nil
}

func main() {
	var queue string
	flag.StringVar(&queue, "q", "", "name of the queue for the worker")
	flag.Parse()

	if err := run(os.Getenv, queue); err != nil {
		slog.Error("Error starting the worker", "error", err)
		os.Exit(1)
	}
}
