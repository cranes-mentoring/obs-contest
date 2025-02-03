package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"syscall"

	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/handler"
	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/logging"
	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/repository"
	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/tracing"
	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/usecase"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// todo: move to cfg
const (
	mongoURI    = "mongodb://mongo1:27017,mongo2:27017,mongo3:27017/?replicaSet=rs0"
	user        = "root"
	pass        = "example"
	dbName      = "project_one"
	timeout     = 30 * time.Second
	maxPoolSize = 10
	minPoolSize = 1
)

func main() {
	ctx := context.Background()

	logging.SetupLogger()
	logger := logging.Logger
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatal("Failed to sync logger:", err)
		}
	}(logger)

	shutdownTracer := tracing.InitTracer(ctx)
	defer shutdownTracer()

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetAuth(options.Credential{Username: user, Password: pass}).
		SetMaxPoolSize(maxPoolSize).
		SetMinPoolSize(minPoolSize).
		SetConnectTimeout(timeout).
		SetServerSelectionTimeout(timeout).
		SetSocketTimeout(timeout)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}

	mongoCtx, mongoCancel := context.WithTimeout(ctx, 10*time.Second)
	defer mongoCancel()

	if err := client.Ping(mongoCtx, nil); err != nil {
		logger.Fatal("Failed to ping MongoDB", zap.String("error", err.Error()))
	}
	logger.Info("Connected to MongoDB!")

	db := client.Database(dbName)

	repo := repository.NewPurchaseRepository(db, logger)
	uc := usecase.NewPurchaseService(repo, logger)
	h := handler.NewPurchaseHandler(uc)

	server := &http.Server{
		Addr: ":38080",
	}

	http.HandleFunc("/api/v1/purchases", h.HandlePurchase)

	go func() {
		logger.Info("Server running.")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("HTTP server ListenAndServe error", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	logger.Info("Shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("HTTP server Shutdown error", zap.Error(err))
	}

	logger.Info("Closing MongoDB connection...")
	if err := client.Disconnect(shutdownCtx); err != nil {
		logger.Fatal("MongoDB Disconnect error", zap.Error(err))
	}

	shutdownTracer()

	logger.Info("Server shut down gracefully")
}
