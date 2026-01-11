package handler

import (
	"clothing-store/internal/service"
	"encoding/json"
	"net/http"
)

type PaymentHandler struct {
	PaymentService *service.PaymentService
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderID  int64  `json:"order_id"`
		Provider string `json:"provider"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	url, err := h.PaymentService.CreatePayment(req.OrderID, req.Provider)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"payment_url": url})
}

// Webhook handles callbacks from payment gateways
func (h *PaymentHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderID       int64  `json:"order_id"`
		TransactionID string `json:"transaction_id"`
		Status        string `json:"status"` // success, failed
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.PaymentService.ConfirmPayment(req.OrderID, req.TransactionID, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}

func (h *PaymentHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	// Basic implementation just to avoid unused method error if I were to add it later?
	// Actually user didn't ask for history endpoint yet, but good for debugging.
	// I'll skip it for now to keep `main.go` simple, or add it if needed.
}
