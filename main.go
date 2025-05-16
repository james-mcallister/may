package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/james-mcallister/may/domain"
)

//go:embed frontend/dist/*
var static embed.FS

func Logger(log *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL)
			next.ServeHTTP(w, r)
		})
	}
}

func home(arg string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Home Page - " + arg + "</h1>"))
	})
}

func main() {
	d := domain.NewDomain()
	d.Init()

	fsys, err := fs.Sub(static, "frontend/dist")
	if err != nil {
		panic(err)
	}
	staticFiles := http.FileServerFS(fsys)

	logger := log.New(os.Stdout, "may:", log.LstdFlags|log.Lshortfile)

	mux := http.NewServeMux()
	middlewareLog := Logger(logger)

	mux.Handle("GET /", middlewareLog(staticFiles))

	// TODO: configurable options
	srv := &http.Server{
		Addr:         "127.0.0.1:54321",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      mux,
	}

	go func() {
		logger.Printf("Starting server on port 54321...")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("server fatal error: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("server shutdown failure: %+v", err)
	}
	logger.Printf("graceful shutdown complete")
	os.Exit(0)
}
