package handler

import (
	"clothing-store/internal/model"
	"clothing-store/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Create a new product (Admin only)
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product body model.Product true "Product Data"
// @Success      201  {object}  model.Product
// @Failure      400  {string}  string "Invalid request"
// @Failure      401  {string}  string "Unauthorized"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if product.Name == "" || product.Price <= 0 {
		http.Error(w, "Name and price are required", http.StatusBadRequest)
		return
	}

	err = h.ProductService.CreateProduct(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// GetProducts godoc
// @Summary      List products
// @Description  Get a list of products with pagination
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number" default(1)
// @Param        limit  query     int  false  "Items per page" default(10)
// @Success      200  {array}   model.Product
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /products [get]
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10

	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscan(p, &page)
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscan(l, &limit)
	}

	products, err := h.ProductService.GetProducts(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.ProductService.GetProductByID(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product model.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.ProductService.UpdateProduct(id, &product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Delete a product by ID (Admin only)
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Product ID"
// @Success      200  {string}  string "Product deleted successfully"
// @Failure      400  {string}  string "Invalid ID"
// @Failure      401  {string}  string "Unauthorized"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.ProductService.DeleteProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product deleted successfully"))
}
