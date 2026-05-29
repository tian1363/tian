package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"note/internal/config"
	"note/internal/db"
	"note/internal/note"
)

func main() {
	cfg := config.Load()

	database, err := db.OpenMySQL(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("open mysql failed: %v", err)
	}
	defer database.Close()

	if err := db.EnsureSchema(database); err != nil {
		log.Fatalf("ensure schema failed: %v", err)
	}

	repo := note.NewRepository(database)
	h := note.NewHandler(repo)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	mux.Handle("GET /", http.FileServer(http.Dir("web")))

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown failed: %v", err)
		}
	}()

	log.Printf("note server listening on %s", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("http server failed: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
