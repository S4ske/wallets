package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	closer2 "wallets/internal/closer"
	config2 "wallets/internal/config"
	"wallets/internal/http-server/handlers"
	"wallets/internal/repositories/postgres"
	"wallets/internal/service"
)

func main() {
	cfg := config2.LoadConfig()

	db, err := sql.Open("postgres", cfg.PostgresURL())
	if err != nil || db == nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	repo := postgres.NewPostgresWalletRepository(db)

	svc := service.NewWalletService(repo, repo, repo)

	h := handlers.NewWalletHandler(svc)

	apiRouter := chi.NewRouter()
	h.RegisterRoutes(apiRouter)

	appRouter := chi.NewRouter()
	appRouter.Mount("/api/v1", apiRouter)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{
		Addr:    cfg.Address,
		Handler: appRouter,
	}

	closer := &closer2.Closer{}
	closer.Add(srv.Shutdown)
	closer.Add(repo.Shutdown)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("server stopping")
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := closer.Close(ctx); err != nil {
		fmt.Printf("%s", err.Error())
	} else {
		fmt.Printf("gracegully shutdowned")
	}
}
