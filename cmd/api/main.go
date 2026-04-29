// Command api starts the HTTP API or, with -seed, runs database seeding.
// Seed mode changes database contents; use only where process identity and DB access are trusted (ops policy).
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"ndinhbang/go-template/database/seeders"
	"ndinhbang/go-template/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, os.Args[1:]); err != nil {
		slog.ErrorContext(ctx, "[main] failed to start", "error", err)
		// os.Exit skips defers — release NotifyContext explicitly so signal hooks are tidy.
		stop()
		os.Exit(1)
	}
}

// run executes the API server or database seeding based on args (typically os.Args[1:]).
// ctx should be the same context used for signal cancellation (e.g. from [signal.NotifyContext]).
// When seeding is enabled, callers must enforce trust boundaries outside this package (deployment ACL, secrets).
func run(ctx context.Context, args []string) error {
	seed, err := parseArgs(args, os.Stderr)
	if err != nil {
		return fmt.Errorf("[main] parse flags: %w", err)
	}
	if seed {
		return runSeeder(ctx)
	}
	return runApi(ctx)
}

// parseArgs parses CLI flags from args. usageOut controls where flag errors write (e.g. [os.Stderr] in production, [io.Discard] in tests).
func parseArgs(args []string, usageOut io.Writer) (seed bool, err error) {
	if usageOut == nil {
		usageOut = os.Stderr
	}
	fs := flag.NewFlagSet("api", flag.ContinueOnError)
	fs.SetOutput(usageOut)
	v := fs.Bool("seed", false, "run database seeding (trusted environments only)")
	if err := fs.Parse(args); err != nil {
		return false, err
	}
	return *v, nil
}

// runApi builds the app via DI and runs until ctx is done.
// Called from [run] and from tests in this package.
func runApi(ctx context.Context) error {
	application, err := app.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("[main] initialize app: %w", err)
	}
	if err := application.Run(ctx); err != nil {
		return fmt.Errorf("[main] run app: %w", err)
	}
	return nil
}

// runSeeder initializes the seeder and runs it.
func runSeeder(ctx context.Context) error {
	seeder, err := seeders.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("[main] initialize seeder: %w", err)
	}
	if err := seeder.Run(ctx); err != nil {
		return fmt.Errorf("[main] seed the database: %w", err)
	}
	slog.Info("[main] database seeded successfully")
	return nil
}
