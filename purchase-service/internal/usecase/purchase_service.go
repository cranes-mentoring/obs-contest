package usecase

import (
	"context"
	"time"

	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/entity"
	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/model"
	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

// purchaseService manages purchase operations using a PurchaseRepository.
type purchaseService struct {
	repo   repository.PurchaseRepository
	logger *zap.Logger
}

// NewPurchaseService initializes a PurchaseUseCase with the provided PurchaseRepository implementation.
func NewPurchaseService(repo repository.PurchaseRepository, logger *zap.Logger) PurchaseUseCase {
	return &purchaseService{repo: repo, logger: logger}
}

// ProcessPurchase handles a purchase request by validating input, storing data in the repository, and returning a response.
func (s *purchaseService) ProcessPurchase(ctx context.Context, request model.PurchaseRequest) (model.PurchaseResponse, error) {
	ctx, span := s.handleTracing(ctx, request.UserID)
	defer span.End()

	purchase := entity.Purchase{
		ID:             uuid.New(),
		UserID:         request.UserID,
		Username:       request.Username,
		Amount:         request.Amount,
		Currency:       request.Currency,
		PaymentMethod:  request.PaymentMethod,
		CardHolderName: request.CardHolderName,
		CardNumber:     request.CardNumber,
		CardExpiry:     request.CardExpiry,
		CardCVC:        request.CardCVC,
		BillingAddress: request.BillingAddress,
		TransactionAt:  request.TransactionAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.SavePurchase(ctx, purchase); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to save purchase")

		return model.PurchaseResponse{}, err
	}

	span.SetStatus(codes.Ok, "Purchase saved successfully")

	response := model.PurchaseResponse{
		TransactionID: purchase.ID.String(),
		Status:        "SUCCESS",
		Message:       "Purchase completed successfully",
	}

	return response, nil
}

func (r *purchaseService) handleTracing(ctx context.Context, userID uuid.UUID) (context.Context, trace.Span) {
	tracer := otel.Tracer("purchase-service")

	ctx, span := tracer.Start(ctx, "ProcessPurchase", trace.WithAttributes(
		attribute.String("user_id", userID.String()),
	))

	return ctx, span
}
