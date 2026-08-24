package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tauan/form-builder/backend/internal/app"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return usage()
	}

	command := args[0]
	cfg, err := app.LoadConfig(args[1:])
	if err != nil {
		return err
	}

	switch command {
	case "run":
		return runServer(cfg)
	case "migrate":
		return migrate(cfg)
	case "openapi":
		return printOpenAPI(cfg)
	default:
		return usage()
	}
}

func runServer(cfg app.Config) error {
	db, err := app.OpenDB(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := app.RunMigrations(db); err != nil {
		return err
	}

	handler, _ := app.NewHandler(cfg, app.NewStore(db))
	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Println("backend listening")
		errCh <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	}
}

func migrate(cfg app.Config) error {
	db, err := app.OpenDB(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := app.RunMigrations(db); err != nil {
		return err
	}
	fmt.Println("migrations applied")
	return nil
}

func printOpenAPI(cfg app.Config) error {
	_, api := app.NewHandler(cfg, nil)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(api.OpenAPI())
}

func usage() error {
	return fmt.Errorf("usage: server <run|migrate|openapi> [--config=./config.yaml] [--address=localhost:8080]")
}
