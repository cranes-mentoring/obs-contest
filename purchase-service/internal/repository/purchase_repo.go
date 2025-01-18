package repository

import (
	"context"
	"time"

	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PurchaseRepository interface {
	SavePurchase(purchase entity.Purchase) error
}

type mongoPurchaseRepo struct {
	collection *mongo.Collection
}

func NewPurchaseRepository(db *mongo.Database) PurchaseRepository {
	return &mongoPurchaseRepo{
		collection: db.Collection("purchases"),
	}
}

func (r *mongoPurchaseRepo) SavePurchase(purchase entity.Purchase) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, bson.M{
		"_id":              purchase.ID,
		"user_id":          purchase.UserID,
		"amount":           purchase.Amount,
		"currency":         purchase.Currency,
		"payment_method":   purchase.PaymentMethod,
		"card_holder_name": purchase.CardHolderName,
		"card_number":      purchase.CardNumber,
		"card_expiry":      purchase.CardExpiry,
		"card_cvc":         purchase.CardCVC,
		"billing_address":  purchase.BillingAddress,
		"transaction_at":   purchase.TransactionAt,
		"created_at":       purchase.CreatedAt,
		"updated_at":       purchase.UpdatedAt,
	})

	return err
}
