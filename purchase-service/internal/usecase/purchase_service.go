package usecase

import (
	"errors"
	"time"

	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/entity"
	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/model"
	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/repository"
	"go.uber.org/zap"

	"github.com/google/uuid"
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
func (s *purchaseService) ProcessPurchase(request model.PurchaseRequest) (model.PurchaseResponse, error) {
	if request.Amount <= 0 {
		return model.PurchaseResponse{}, errors.New("invalid amount")
	}
	if request.Currency == "" {
		return model.PurchaseResponse{}, errors.New("currency is required")
	}

	purchase := entity.Purchase{
		ID:             uuid.New(),
		UserID:         request.UserID,
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

	if err := s.repo.SavePurchase(purchase); err != nil {
		return model.PurchaseResponse{}, err
	}

	response := model.PurchaseResponse{
		TransactionID: purchase.ID.String(),
		Status:        "SUCCESS",
		Message:       "Purchase completed successfully",
	}

	return response, nil
}
