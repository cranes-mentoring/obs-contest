package usecase

import (
	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/model"
)

type PurchaseUseCase interface {
	ProcessPurchase(request model.PurchaseRequest) (model.PurchaseResponse, error)
}
