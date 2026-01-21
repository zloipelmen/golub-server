package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"messenger/internal/auth"
	"messenger/internal/config"
	"messenger/internal/httpx"
	"messenger/internal/storage"
	"messenger/internal/ws"
)

func main() {
	cfg := config.FromEnv()

	db, err := storage.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	store := storage.NewStore(db)
	authSvc := auth.NewService(store)

	hub := ws.NewHub()
	go hub.Run()

	router := httpx.NewRouter(hub, authSvc, store)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("shutdown complete")
}
