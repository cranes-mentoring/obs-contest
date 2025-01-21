package usecase

import (
	"context"

	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/model"
)

type PurchaseUseCase interface {
	ProcessPurchase(ctx context.Context, request model.PurchaseRequest) (model.PurchaseResponse, error)
}
