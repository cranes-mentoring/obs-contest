package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/handler"
	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/repository"
	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetAuth(options.Credential{Username: "example_user", Password: "example_password"})

	clientOptions.SetMaxPoolSize(10)
	clientOptions.SetMinPoolSize(1)
	clientOptions.SetConnectTimeout(10 * time.Second)
	clientOptions.SetServerSelectionTimeout(10 * time.Second)
	clientOptions.SetSocketTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	log.Println("Connected to MongoDB!")

	db := client.Database("purchase_service")

	repo := repository.NewPurchaseRepository(db)
	uc := usecase.NewPurchaseService(repo)
	h := handler.NewPurchaseHandler(uc)

	http.HandleFunc("/api/v1/purchases", h.HandlePurchase)

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
