package repository

import (
	"context"

	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PurchaseRepository interface {
	SavePurchase(ctx context.Context, purchase entity.Purchase) error
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

func (r *mongoPurchaseRepo) SavePurchase(ctx context.Context, purchase entity.Purchase) error {
	span := r.handleTracing(ctx, purchase)
	defer span.End()

	_, err := r.collection.InsertOne(ctx, bson.M{
		"_id":              purchase.ID,
		"user_id":          purchase.UserID,
		"username":         purchase.Username,
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
		"trace_id":         span.SpanContext().TraceID().String(),
	})

	return err
}

func (r *mongoPurchaseRepo) handleTracing(ctx context.Context, purchase entity.Purchase) trace.Span {
	tracer := otel.Tracer("purchase-service")

	_, span := tracer.Start(ctx, "SavePurchase", trace.WithAttributes(
		attribute.String("db.system", "mongodb"),
		attribute.String("db.collection", r.collection.Name()),
		attribute.String("purchase.user_id", purchase.UserID.String()),
	))

	return span
}
