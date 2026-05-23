package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikitagrgv/shorter/internal/config"
	deliveryHttp "github.com/nikitagrgv/shorter/internal/delivery/http"
	"github.com/nikitagrgv/shorter/internal/infrastructure/id"
	"github.com/nikitagrgv/shorter/internal/infrastructure/postgres"
	"github.com/nikitagrgv/shorter/internal/infrastructure/shortlink"
	infrTime "github.com/nikitagrgv/shorter/internal/infrastructure/time"
	"github.com/nikitagrgv/shorter/internal/infrastructure/token"
	"github.com/nikitagrgv/shorter/internal/infrastructure/token_hasher"
	"github.com/nikitagrgv/shorter/internal/usecase"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.Database,
	)

	fmt.Println("Database Connecting...")
	pool, err := postgres.NewPool(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Error: %v\nConnection string: %s", err, connStr)
	}

	fmt.Println("Database Connected")

	defer func() {
		fmt.Println("Database Disconnecting...")
		pool.Close()
		fmt.Println("Database Disconnected")
	}()

	mux := http.NewServeMux()

	staticFS, err := fs.Sub(deliveryHttp.Assets, "static")
	if err != nil {
		log.Fatal(err)
	}

	staticHandler := http.FileServer(http.FS(staticFS))
	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	tmpl, err := template.ParseFS(deliveryHttp.Assets, "templates/*.html")
	if err != nil {
		log.Fatalf("Cannot parse templates: %v", err)
	}

	timeProvider := &infrTime.RealTimeProvider{}
	repo := postgres.NewLinkRepoPostgres(pool)
	idGen, err := id.NewSnowflakeIdGenerator(cfg.App.NodeId)
	shortGen := shortlink.NewHashidsLinkEncoder()
	tokenGen := &token.UuidTokenGenerator{}
	tokenHasher := &token_hasher.TokenHasher{}
	if err != nil {
		log.Fatalf("Cannot init snowflake: %v", err)
	}

	linkUsecase := usecase.NewLinkUsecase(timeProvider, repo, idGen, shortGen, tokenGen, tokenHasher)

	handler := deliveryHttp.NewLinkHandler(tmpl, linkUsecase, cfg.App.BaseUrl)

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		handler.ShowHello(w, r)
	})

	mux.HandleFunc("GET /{code}", func(w http.ResponseWriter, r *http.Request) {
		handler.RedirectShortLink(w, r)
	})

	mux.HandleFunc("GET /links/{token}/manage", func(w http.ResponseWriter, r *http.Request) {
		handler.ShowResult(w, r)
	})

	mux.HandleFunc("POST /api/v1/links", func(w http.ResponseWriter, r *http.Request) {
		handler.PostLink(w, r)
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.ShorterPort),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Println("Server is running on addr ", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %s\n", err)
		}
	}()

	<-stop
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}
