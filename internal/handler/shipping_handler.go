package handler

import (
	"clothing-store/internal/model"
	"clothing-store/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type UpdateShippingStatusRequest struct {
	ShippingID int64  `json:"shipping_id"`
	Status     string `json:"status" example:"shipping"`
	Location   string `json:"location" example:"Warehouse A"`
	Note       string `json:"note" example:"Picked up"`
}

type UpdateTrackingCodeRequest struct {
	TrackingCode string `json:"tracking_code" example:"GHN123456"`
}

type ShippingHandler struct {
	ShippingService *service.ShippingService
}

// CreateShipping godoc
// @Summary      Create shipping info
// @Description  Create shipping information for an order
// @Tags         shippings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        shipping body model.Shipping true "Shipping Data"
// @Success      201  {object}  model.Shipping
// @Failure      400  {string}  string "Invalid request"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /shippings [post]
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

// GetShippingByOrder godoc
// @Summary      Get shipping by order ID
// @Description  Get shipping information by order ID
// @Tags         shippings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_id query int true "Order ID"
// @Success      200  {object}  model.Shipping
// @Failure      400  {string}  string "Invalid order ID"
// @Failure      404  {string}  string "Shipping not found"
// @Router       /shippings [get]
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

// TrackShipping godoc
// @Summary      Track shipping
// @Description  Track shipping status by tracking code
// @Tags         shippings
// @Accept       json
// @Produce      json
// @Param        tracking_code query string true "Tracking Code"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {string}  string "Tracking code is required"
// @Failure      404  {string}  string "Tracking code not found"
// @Router       /shippings/track [get]
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

// UpdateShippingStatus godoc
// @Summary      Update shipping status
// @Description  Update status of a shipping (Admin/System)
// @Tags         shippings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdateShippingStatusRequest true "Status Update"
// @Success      200  {object}  map[string]string
// @Failure      400  {string}  string "Invalid request"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /shippings/status [put]
func (h *ShippingHandler) UpdateShippingStatus(w http.ResponseWriter, r *http.Request) {
	var req UpdateShippingStatusRequest

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

// GetShippingHistory godoc
// @Summary      Get shipping history
// @Description  Get history of a shipping
// @Tags         shippings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        shipping_id query int true "Shipping ID"
// @Success      200  {array}   model.ShippingHistory
// @Failure      400  {string}  string "Invalid shipping ID"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /shippings/history [get]
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

// UpdateTrackingCode godoc
// @Summary      Update tracking code
// @Description  Update tracking code for a shipping
// @Tags         shippings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Shipping ID"
// @Param        request body UpdateTrackingCodeRequest true "Tracking Code"
// @Success      200  {string}  string "Tracking code updated successfully"
// @Failure      400  {string}  string "Invalid request"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /shippings/{id}/tracking [put]
func (h *ShippingHandler) UpdateTrackingCode(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/shippings/")
	idStr = strings.TrimSuffix(idStr, "/tracking")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid shipping ID", http.StatusBadRequest)
		return
	}

	var req UpdateTrackingCodeRequest

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
