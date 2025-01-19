package repository

import (
	"context"
	"time"

	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PurchaseRepository interface {
	SavePurchase(purchase entity.Purchase) error
}

type mongoPurchaseRepo struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewPurchaseRepository(db *mongo.Database, logger *zap.Logger) PurchaseRepository {
	return &mongoPurchaseRepo{
		collection: db.Collection("purchases"),
		logger:     logger,
	}
}

func (r *mongoPurchaseRepo) SavePurchase(purchase entity.Purchase) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	span := r.handleTracing(ctx, purchase)
	defer span.End()

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

func (r *mongoPurchaseRepo) handleTracing(ctx context.Context, purchase entity.Purchase) trace.Span {
	tracer := otel.Tracer("purchase-repository")

	_, span := tracer.Start(ctx, "SavePurchase", trace.WithAttributes(
		attribute.String("db.system", "mongodb"),
		attribute.String("db.name", r.collection.Database().Name()),
		attribute.String("db.collection", r.collection.Name()),
		attribute.String("purchase.user_id", purchase.UserID.String()),
		attribute.Float64("purchase.amount", purchase.Amount),
		attribute.String("purchase.currency", purchase.Currency),
	))

	return span
}
