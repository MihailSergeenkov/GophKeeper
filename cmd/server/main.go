package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/MihailSergeenkov/GophKeeper/internal/logger"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/crypt"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/handlers"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/routes"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/s3"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/services"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
)

const (
	timeoutServerShutdown = time.Second * 5
	timeoutShutdown       = time.Second * 10
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

// WebServer интерфейс к веб серверу.
type WebServer interface {
	ListenAndServeTLS(certFile string, keyFile string) error
	ListenAndServe() error
}

func main() {
	ctx := context.Background()
	withFlags := true
	if err := run(ctx, withFlags); err != nil {
		log.Fatal(err)
	}
	log.Println("bye-bye")
}

func run(baseCtx context.Context, withFlags bool) error {
	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)

	ctx, cancelCtx := signal.NotifyContext(baseCtx, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancelCtx()

	g, ctx := errgroup.WithContext(ctx)

	context.AfterFunc(ctx, func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), timeoutShutdown)
		defer cancelCtx()

		<-ctx.Done()
		log.Fatal("failed to gracefully shutdown the service")
	})

	c, err := config.Setup(withFlags)
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	l, err := logger.NewLogger(c.LogLevel)
	if err != nil {
		return fmt.Errorf("logger error: %w", err)
	}

	l.Info("Running server on", zap.String("addr", c.RunAddr))
	l.Info("S3 server on", zap.String("addr", c.S3.Endpoint))

	store, err := storage.NewStorage(ctx, l, c.DatabaseURI)
	if err != nil {
		return fmt.Errorf("storage error: %w", err)
	}

	g.Go(func() error {
		defer log.Print("closed DB")

		<-ctx.Done()

		if err := store.Close(); err != nil {
			l.Error("failed to close db connection", zap.Error(err))
		}
		return nil
	})

	fs, err := s3.NewClient(ctx, &c.S3)
	if err != nil {
		return fmt.Errorf("s3 error: %w", err)
	}

	cr, err := crypt.NewCrypt(c)
	if err != nil {
		return fmt.Errorf("crypt error: %w", err)
	}

	s := services.NewServices(store, fs, cr, c)
	h := handlers.NewHandlers(s, l)
	r := routes.NewRouter(h, c, l, store)

	srv := configureServer(r, c.EnableHTTPS, c.RunAddr)

	g.Go(func() error {
		defer func() {
			errRec := recover()
			if errRec != nil {
				err = fmt.Errorf("a panic occurred: %v", errRec)
				l.Error("failed", zap.Error(err))
			}
		}()
		if err := runServer(srv, c.EnableHTTPS); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				return fmt.Errorf("HTTP server has encoutenred an error: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		defer log.Print("server has been shutdown")
		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()
		if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
			log.Printf("an error occurred during server shutdown: %v", err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("some errorgroup error: %w", err)
	}

	return nil
}

func configureServer(r chi.Router, enableHTTPS bool, runAddr string) *http.Server {
	if enableHTTPS {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("/tmp/certs"),
			HostPolicy: autocert.HostWhitelist("mynetwork.keenetic.link"),
		}
		server := &http.Server{
			Addr:    runAddr,
			Handler: r,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
				MinVersion:     tls.VersionTLS13,
			},
		}

		return server
	}

	return &http.Server{
		Addr:    runAddr,
		Handler: r,
	}
}

func runServer(srv WebServer, enableHTTPS bool) error {
	var err error

	if enableHTTPS {
		err = srv.ListenAndServeTLS("", "")
	} else {
		err = srv.ListenAndServe()
	}
	if err != nil {
		return fmt.Errorf("listen and server has failed: %w", err)
	}

	return nil
}
