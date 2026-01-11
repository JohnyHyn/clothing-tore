package handler

import (
	"clothing-store/internal/model"
	"clothing-store/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ShippingHandler struct {
	ShippingService *service.ShippingService
}

// CreateShipping tạo thông tin vận chuyển cho order
func (h *ShippingHandler) CreateShipping(w http.ResponseWriter, r *http.Request) {
	var shipping model.Shipping

	if err := json.NewDecoder(r.Body).Decode(&shipping); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validation
	if shipping.OrderID == 0 {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	// Calculate shipping fee if not provided
	if shipping.Fee == 0 {
		fee, err := h.ShippingService.CalculateShippingFee(shipping.Method, shipping.City)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shipping.Fee = fee
	}

	// Set default status
	if shipping.Status == "" {
		shipping.Status = "pending"
	}

	err := h.ShippingService.CreateShipping(&shipping)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Shipping created successfully",
		"data":    shipping,
	})
}

// GetShippingByOrder lấy thông tin shipping theo order ID
func (h *ShippingHandler) GetShippingByOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("order_id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	shipping, err := h.ShippingService.GetShippingByOrderID(orderID)
	if err != nil {
		http.Error(w, "Shipping not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": shipping,
	})
}

// TrackShipping theo dõi đơn hàng
func (h *ShippingHandler) TrackShipping(w http.ResponseWriter, r *http.Request) {
	trackingCode := r.URL.Query().Get("tracking_code")
	if trackingCode == "" {
		http.Error(w, "Tracking code is required", http.StatusBadRequest)
		return
	}

	shipping, history, err := h.ShippingService.TrackShipping(trackingCode)
	if err != nil {
		http.Error(w, "Tracking code not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"shipping": shipping,
		"history":  history,
	})
}

// UpdateShippingStatus cập nhật trạng thái vận chuyển
func (h *ShippingHandler) UpdateShippingStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ShippingID int64  `json:"shipping_id"`
		Status     string `json:"status"`
		Location   string `json:"location"`
		Note       string `json:"note"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.ShippingID == 0 || req.Status == "" {
		http.Error(w, "Shipping ID and status are required", http.StatusBadRequest)
		return
	}

	err := h.ShippingService.UpdateShippingStatus(req.ShippingID, req.Status, req.Location, req.Note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Shipping status updated successfully",
	})
}

// GetShippingHistory lấy lịch sử vận chuyển
func (h *ShippingHandler) GetShippingHistory(w http.ResponseWriter, r *http.Request) {
	shippingIDStr := r.URL.Query().Get("shipping_id")
	shippingID, err := strconv.ParseInt(shippingIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid shipping ID", http.StatusBadRequest)
		return
	}

	history, err := h.ShippingService.GetShippingHistory(shippingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": history,
	})
}

// CalculateShippingFee tính phí ship
func (h *ShippingHandler) CalculateShippingFee(w http.ResponseWriter, r *http.Request) {
	method := r.URL.Query().Get("method")
	city := r.URL.Query().Get("city")

	if method == "" || city == "" {
		http.Error(w, "Method and city are required", http.StatusBadRequest)
		return
	}

	fee, err := h.ShippingService.CalculateShippingFee(method, city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"method": method,
		"city":   city,
		"fee":    fee,
	})
}

// UpdateTrackingCode cập nhật mã vận đơn
func (h *ShippingHandler) UpdateTrackingCode(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/shippings/")
	idStr = strings.TrimSuffix(idStr, "/tracking")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid shipping ID", http.StatusBadRequest)
		return
	}

	var req struct {
		TrackingCode string `json:"tracking_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = h.ShippingService.UpdateTrackingCode(id, req.TrackingCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Tracking code updated successfully"))
}
