package handler

import (
	"clothing-store/internal/model"
	"clothing-store/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type VoucherHandler struct {
	VoucherService *service.VoucherService
}

// CreateVoucher tạo voucher mới (admin only)
func (h *VoucherHandler) CreateVoucher(w http.ResponseWriter, r *http.Request) {
	var voucher model.Voucher

	if err := json.NewDecoder(r.Body).Decode(&voucher); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validation
	if voucher.Code == "" || voucher.Type == "" {
		http.Error(w, "Code and type are required", http.StatusBadRequest)
		return
	}

	if voucher.Type != "percentage" && voucher.Type != "fixed" {
		http.Error(w, "Type must be 'percentage' or 'fixed'", http.StatusBadRequest)
		return
	}

	err := h.VoucherService.CreateVoucher(&voucher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Voucher created successfully",
		"data":    voucher,
	})
}

// ValidateVoucher kiểm tra voucher có hợp lệ không
func (h *VoucherHandler) ValidateVoucher(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code        string  `json:"code"`
		OrderAmount float64 `json:"order_amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	voucher, discount, err := h.VoucherService.ValidateVoucher(req.Code, req.OrderAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":    true,
		"voucher":  voucher,
		"discount": discount,
		"message":  "Voucher is valid",
	})
}

// ListVouchers lấy danh sách vouchers
func (h *VoucherHandler) ListVouchers(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	activeOnlyStr := r.URL.Query().Get("active_only")

	page := 1
	limit := 10
	activeOnly := false

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

	if activeOnlyStr == "true" || activeOnlyStr == "1" {
		activeOnly = true
	}

	vouchers, total, err := h.VoucherService.ListVouchers(page, limit, activeOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (total + limit - 1) / limit

	response := map[string]interface{}{
		"data": vouchers,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetVoucher lấy thông tin voucher theo code
func (h *VoucherHandler) GetVoucher(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/vouchers/")

	voucher, err := h.VoucherService.GetVoucherByCode(code)
	if err != nil {
		http.Error(w, "Voucher not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": voucher,
	})
}

// UpdateVoucher cập nhật voucher (admin only)
func (h *VoucherHandler) UpdateVoucher(w http.ResponseWriter, r *http.Request) {
	var voucher model.Voucher

	if err := json.NewDecoder(r.Body).Decode(&voucher); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.VoucherService.UpdateVoucher(&voucher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Voucher updated successfully",
		"data":    voucher,
	})
}

// DeleteVoucher xóa voucher (admin only)
func (h *VoucherHandler) DeleteVoucher(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/vouchers/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid voucher ID", http.StatusBadRequest)
		return
	}

	err = h.VoucherService.DeleteVoucher(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Voucher deleted successfully"))
}
