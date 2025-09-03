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
	config2 "wallets/internal/config"
	"wallets/internal/http-server/handlers"
	"wallets/internal/repositories/postgres"
	"wallets/internal/service"
)

func main() {
	cfg := config2.LoadConfig()
	db, err := sql.Open("postgres", cfg.PostgresURL())
	if err != nil {
		log.Fatalf("failed to connect db: %s", err.Error())
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

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("server stopping")
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("failed to stop server: %s", err.Error())
	} else {
		fmt.Printf("gracegully shutdowned")
	}
}
