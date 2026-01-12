package handler

import (
	"clothing-store/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type CartHandler struct {
	CartService *service.CartService
}

type AddToCartRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// GetCart godoc
// @Summary      Get user's cart
// @Description  Get all items in the current user's cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   model.CartItem
// @Failure      401  {string}  string "Unauthorized"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /cart [get]
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64) // From JWT Middleware
	items, err := h.CartService.GetCart(int64(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// AddToCart godoc
// @Summary      Add item to cart
// @Description  Add a product to the user's cart or increase quantity
// @Tags         cart
// @Accept       json
// @Produce      plain
// @Security     BearerAuth
// @Param        request body AddToCartRequest true "Product and Quantity"
// @Success      201  {string}  string "Added to cart"
// @Failure      400  {string}  string "Invalid request"
// @Failure      401  {string}  string "Unauthorized"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /cart [post]
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)
	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}
	err := h.CartService.AddToCart(int64(userID), req.ProductID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Added to cart"))
}

// UpdateQuantity godoc
// @Summary      Update cart item quantity
// @Description  Change the quantity of an item in the cart
// @Tags         cart
// @Accept       json
// @Produce      plain
// @Security     BearerAuth
// @Param        id   path      int  true  "Cart Item ID"
// @Param        request body object true "New quantity"
// @Success      200  {string}  string "Quantity updated"
// @Failure      401  {string}  string "Unauthorized"
// @Router       /cart/{id} [put]
func (h *CartHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)
	idStr := strings.TrimPrefix(r.URL.Path, "/cart/")
	cartItemID, _ := strconv.ParseInt(idStr, 10, 64)

	var req struct {
		Quantity int `json:"quantity"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Quantity <= 0 {
		h.CartService.RemoveItem(int64(userID), cartItemID)
	} else {
		h.CartService.UpdateQuantity(int64(userID), cartItemID, req.Quantity)
	}
	w.Write([]byte("Quantity updated"))
}

// RemoveItem godoc
// @Summary      Remove item from cart
// @Description  Remove a single item from the cart by its ID
// @Tags         cart
// @Accept       json
// @Produce      plain
// @Security     BearerAuth
// @Param        id   path      int  true  "Cart Item ID"
// @Success      200  {string}  string "Item removed"
// @Failure      401  {string}  string "Unauthorized"
// @Router       /cart/{id} [delete]
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)
	idStr := strings.TrimPrefix(r.URL.Path, "/cart/")
	cartItemID, _ := strconv.ParseInt(idStr, 10, 64)

	err := h.CartService.RemoveItem(int64(userID), cartItemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Item removed"))
}
