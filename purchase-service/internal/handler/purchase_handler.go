package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/model"
	"github.com/cranes-mentoring/obs-contest/purchase-service/internal/usecase"
)

type PurchaseHandler struct {
	UseCase usecase.PurchaseUseCase
}

func NewPurchaseHandler(uc usecase.PurchaseUseCase) *PurchaseHandler {
	return &PurchaseHandler{UseCase: uc}
}

func (h *PurchaseHandler) HandlePurchase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request model.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)

		return
	}

	response, err := h.UseCase.ProcessPurchase(ctx, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
