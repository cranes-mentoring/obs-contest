package main

import (
	"context"
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

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Global context and shutdown setup
	ctx := context.Background()

	// Set up logging
	logging.SetupLogger()
	logger := logging.Logger
	defer logger.Sync()

	// Initialize tracing
	shutdownTracer := tracing.InitTracer(ctx)
	defer shutdownTracer()

	// MongoDB connection setup
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetAuth(options.Credential{Username: "example_user", Password: "example_password"}).
		SetMaxPoolSize(10).
		SetMinPoolSize(1).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second).
		SetSocketTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	mongoCtx, mongoCancel := context.WithTimeout(ctx, 5*time.Second)
	defer mongoCancel()

	// Verify connection to MongoDB
	if err := client.Ping(mongoCtx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("Connected to MongoDB!")

	db := client.Database("project_one")

	// Initialize business logic
	repo := repository.NewPurchaseRepository(db, logger)
	uc := usecase.NewPurchaseService(repo, logger)
	h := handler.NewPurchaseHandler(uc)

	// Create HTTP server
	server := &http.Server{
		Addr: ":8080",
	}

	// HTTP endpoint setup
	http.HandleFunc("/api/v1/purchases", h.HandlePurchase)

	// Start server in a goroutine
	go func() {
		log.Println("Server running at :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Wait for a termination signal
	log.Println("Shutdown signal received")

	// Timeout context for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	log.Println("Shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server Shutdown error: %v", err)
	}

	// Close MongoDB connection
	log.Println("Closing MongoDB connection...")
	if err := client.Disconnect(shutdownCtx); err != nil {
		log.Fatalf("MongoDB Disconnect error: %v", err)
	}

	// Stop tracer if necessary
	shutdownTracer()

	log.Println("Server shut down gracefully")
}
