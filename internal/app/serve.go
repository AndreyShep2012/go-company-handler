package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndreyShep2012/go-company-handler/internal/config"
	"golang.org/x/sync/errgroup"
)

func Serve(config config.Config) {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	initLogger(config.LogLevel)
	fiberServer, api := initFiberServer(config.ApiRoot, config.JWTSecretKey)
	collection := initMongo(mainCtx, config.MongoUri, config.MongoDatabaseName, config.MongoCompaniesCollection, config.ConnectTimeoutSec)
	setupRoutes(config, fiberServer, api, collection)

	g, gCtx := errgroup.WithContext(mainCtx)

	g.Go(func() error {
		return fiberServer.Listen(config.ListenAddr)
	})

	g.Go(func() error {
		<-gCtx.Done()
		fiberServer.Shutdown()
		slog.Info("server shutdown")
		return gCtx.Err()
	})

	if err := g.Wait(); err != nil {
		switch err {
		case http.ErrServerClosed:
		case context.Canceled:
		default:
			slog.Error("error running service: ", "error", err.Error())
		}
	}
}
