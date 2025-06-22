package main

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//go:embed frontend/dist/*
var static embed.FS

//go:embed templates/*
var templ embed.FS

// addition method for template iteration
func add(a, b int) int {
	return a + b
}

func main() {
	d := NewDomain()
	d.Init()

	funcMap := template.FuncMap{
		"add": add,
	}
	d.templates = template.Must(template.New("may-templates").Funcs(funcMap).ParseFS(templ, "templates/*.html"))

	fsys, err := fs.Sub(static, "frontend/dist")
	if err != nil {
		panic(err)
	}
	staticFiles := http.FileServerFS(fsys)

	logger := log.New(os.Stdout, "may:", log.LstdFlags|log.Lshortfile)
	middlewareLog := Logger(logger)

	mux := http.NewServeMux()
	mux.Handle("/", middlewareLog(staticFiles))
	initRoutes(mux, logger, d)

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
