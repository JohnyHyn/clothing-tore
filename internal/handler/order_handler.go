package handler

import (
	"clothing-store/internal/model"
	"clothing-store/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type RefundRequest struct {
	Reason string `json:"reason" example:"Product defective"`
}

type OrderHandler struct {
	OrderService *service.OrderService
}

// CreateOrder godoc
// @Summary      Create a new order
// @Description  Create a new order with items
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order body model.Order true "Order Data"
// @Success      201  {object}  model.Order
// @Failure      400  {string}  string "Invalid request"
// @Failure      401  {string}  string "Unauthorized"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(float64)
	order.UserID = int64(userID)

	if len(order.Items) == 0 {
		http.Error(w, "Order items required", http.StatusBadRequest)
		return
	}

	err := h.OrderService.CreateOrder(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// CancelOrder godoc
// @Summary      Cancel an order
// @Description  Cancel an order by ID
// @Tags         orders
// @Accept       json
// @Produce      plain
// @Security     BearerAuth
// @Param        id   path      int  true  "Order ID"
// @Success      200  {string}  string "Order cancelled successfully"
// @Failure      400  {string}  string "Invalid ID"
// @Failure      401  {string}  string "Unauthorized"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /orders/{id}/cancel [put]
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	idStr = strings.TrimSuffix(idStr, "/cancel")

	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	err = h.OrderService.CancelOrder(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Order cancelled successfully"))
}

// DeleteOrder godoc
// @Summary      Delete an order
// @Description  Delete an order by ID (only pending or cancelled)
// @Tags         orders
// @Accept       json
// @Produce      plain
// @Security     BearerAuth
// @Param        id   path      int  true  "Order ID"
// @Success      200  {string}  string "Order deleted successfully"
// @Failure      400  {string}  string "Invalid ID"
// @Failure      401  {string}  string "Unauthorized"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	err = h.OrderService.DeleteOrder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order deleted successfully"))
}

// PayOrder godoc
// @Summary      Pay an order
// @Description  Mark an order as paid
// @Tags         orders
// @Accept       json
// @Produce      plain
// @Security     BearerAuth
// @Param        id   path      int  true  "Order ID"
// @Success      200  {string}  string "Order paid successfully"
// @Failure      400  {string}  string "Invalid ID"
// @Failure      401  {string}  string "Unauthorized"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /orders/{id}/pay [put]
func (h *OrderHandler) PayOrder(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	idStr = strings.TrimSuffix(idStr, "/pay")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	err = h.OrderService.PayOrder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order paid successfully"))
}

// ApproveOrder godoc
// @Summary      Approve an order
// @Description  Change order status from 'chờ xử lý' to 'đã xử lý' (Staff/Admin only)
// @Tags         orders
// @Accept       json
// @Produce      plain
// @Security     BearerAuth
// @Param        id   path      int  true  "Order ID"
// @Success      200  {string}  string "Order approved successfully"
// @Failure      400  {string}  string "Invalid ID or status"
// @Router       /orders/{id}/approve [put]
func (h *OrderHandler) ApproveOrder(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	idStr = strings.TrimSuffix(idStr, "/approve")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	err := h.OrderService.ApproveOrder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("Order approved successfully"))
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.OrderService.GetOrderByID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": order,
	})
}

// RefundOrder godoc
// @Summary      Refund an order
// @Description  Refund an order with a reason
// @Tags         orders
// @Accept       json
// @Produce      plain
// @Security     BearerAuth
// @Param        id   path      int  true  "Order ID"
// @Param        request body RefundRequest true "Refund Reason"
// @Success      200  {string}  string "Order refunded successfully"
// @Failure      400  {string}  string "Invalid ID or request"
// @Failure      401  {string}  string "Unauthorized"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /orders/{id}/refund [put]
func (h *OrderHandler) RefundOrder(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	idStr = strings.TrimSuffix(idStr, "/refund")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	var req RefundRequest
	json.NewDecoder(r.Body).Decode(&req)

	err = h.OrderService.RefundOrder(id, req.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Order refunded successfully"))
}

// ListOrders handles GET requests to list orders with pagination and filters
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters from query string
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Parse filter parameters
	status := r.URL.Query().Get("status")
	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Get orders from service with filters
	orders, total, err := h.OrderService.ListOrders(page, limit, status, dateFrom, dateTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate pagination metadata
	totalPages := (total + limit - 1) / limit

	// Prepare response
	response := map[string]interface{}{
		"data": orders,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
		"filters": map[string]interface{}{
			"status":    status,
			"date_from": dateFrom,
			"date_to":   dateTo,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListMyOrders handles GET requests to list the current user's orders
// @Summary      List my orders
// @Description  Get a paginated list of orders for the currently logged-in user
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int  false  "Page number" default(1)
// @Param        limit  query     int  false  "Items per page" default(10)
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {string}  string "Unauthorized"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /my/orders [get]
func (h *OrderHandler) ListMyOrders(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("user_id")
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDVal.(float64)

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	orders, total, err := h.OrderService.ListUserOrders(int64(userID), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data": orders,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
