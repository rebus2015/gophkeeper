package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/rebus2015/gophkeeper/internal/logger"
	"github.com/rebus2015/gophkeeper/internal/server"
	"github.com/rebus2015/gophkeeper/internal/server/config"
	"github.com/rebus2015/gophkeeper/internal/storage/db"
	"github.com/rebus2015/gophkeeper/internal/storage/migrations"
)

var (
	buildVersion     = "N/A"
	buildDate        = "N/A"
	buildCommit      = "N/A"
	fileReadTimeout  = 30 * time.Second
	fileWriteTimeout = 30 * time.Second
)

func main() {

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n\n", buildCommit)

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Panicf("Error reading configuration from env variables: %v", err)
		return
	}

	lg := logger.NewConsole(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = migrations.RunMigrations(lg, cfg)
	if err != nil {
		lg.Fatal().Err(err).Msgf("Migrations retuned error")
		return
	}
	storage, err := db.NewStorage(ctx, lg, cfg)
	if err != nil {
		lg.Fatal().Err(err).Msgf("Error creating dbStorage, with conn: %s", cfg.ConnectionString)
		return
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		<-sigChan
		cancel()
	}()

	g, gCtx := errgroup.WithContext(ctx)
	grpcSrv := server.NewRPCServer(*storage, *cfg, lg)
	g.Go(func() error {
		<-gCtx.Done()
		grpcSrv.Shutdown()
		log.Println("GRPC Server shutdown gracefully!")
		return nil
	})

	g.Go(func() error {
		if err := grpcSrv.Run(); err != nil {
			// ошибки запуска Listener
			log.Printf("Error gRPC server Start: %v", err)
			return fmt.Errorf("gRPC server Start error: %w", err)
		}
		return nil
	})

	err = g.Wait()
	if err != nil {
		log.Printf("error: server exited with %v", err)
	}
	fmt.Println("Server Shutdown gracefully")
}
