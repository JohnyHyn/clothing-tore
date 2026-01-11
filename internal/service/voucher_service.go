package service

import (
	"clothing-store/internal/model"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type VoucherService struct {
	DB *sql.DB
}

// CreateVoucher tạo voucher mới
func (s *VoucherService) CreateVoucher(voucher *model.Voucher) error {
	query := `
		INSERT INTO vouchers (code, description, type, value, min_order, max_discount, 
			usage_limit, used_count, start_date, end_date, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?)
	`

	result, err := s.DB.Exec(query,
		strings.ToUpper(voucher.Code),
		voucher.Description,
		voucher.Type,
		voucher.Value,
		voucher.MinOrder,
		voucher.MaxDiscount,
		voucher.UsageLimit,
		voucher.StartDate,
		voucher.EndDate,
		voucher.IsActive,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	voucher.ID = id
	return nil
}

// GetVoucherByCode lấy voucher theo mã
func (s *VoucherService) GetVoucherByCode(code string) (*model.Voucher, error) {
	voucher := &model.Voucher{}

	query := `
		SELECT id, code, description, type, value, min_order, max_discount,
			usage_limit, used_count, start_date, end_date, is_active, created_at
		FROM vouchers
		WHERE code = ?
	`

	err := s.DB.QueryRow(query, strings.ToUpper(code)).Scan(
		&voucher.ID,
		&voucher.Code,
		&voucher.Description,
		&voucher.Type,
		&voucher.Value,
		&voucher.MinOrder,
		&voucher.MaxDiscount,
		&voucher.UsageLimit,
		&voucher.UsedCount,
		&voucher.StartDate,
		&voucher.EndDate,
		&voucher.IsActive,
		&voucher.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return voucher, nil
}

// ValidateVoucher kiểm tra voucher có hợp lệ không
func (s *VoucherService) ValidateVoucher(code string, orderAmount float64) (*model.Voucher, float64, error) {
	voucher, err := s.GetVoucherByCode(code)
	if err != nil {
		return nil, 0, errors.New("voucher not found")
	}

	// Check if active
	if !voucher.IsActive {
		return nil, 0, errors.New("voucher is not active")
	}

	// Check usage limit
	if voucher.UsedCount >= voucher.UsageLimit {
		return nil, 0, errors.New("voucher usage limit reached")
	}

	// Check date range
	now := time.Now()
	if now.Before(voucher.StartDate) {
		return nil, 0, errors.New("voucher not yet valid")
	}
	if now.After(voucher.EndDate) {
		return nil, 0, errors.New("voucher has expired")
	}

	// Check minimum order
	if orderAmount < voucher.MinOrder {
		return nil, 0, errors.New("order amount below minimum")
	}

	// Calculate discount
	var discount float64
	if voucher.Type == "percentage" {
		discount = orderAmount * (voucher.Value / 100)
		if voucher.MaxDiscount > 0 && discount > voucher.MaxDiscount {
			discount = voucher.MaxDiscount
		}
	} else {
		discount = voucher.Value
	}

	return voucher, discount, nil
}

// ApplyVoucher áp dụng voucher cho order
func (s *VoucherService) ApplyVoucher(voucherID, orderID int64, discount float64) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Increment used_count
	_, err = tx.Exec("UPDATE vouchers SET used_count = used_count + 1 WHERE id = ?", voucherID)
	if err != nil {
		return err
	}

	// Record usage
	_, err = tx.Exec(`
		INSERT INTO voucher_usage (voucher_id, order_id, discount)
		VALUES (?, ?, ?)
	`, voucherID, orderID, discount)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// ListVouchers lấy danh sách vouchers
func (s *VoucherService) ListVouchers(page, limit int, activeOnly bool) ([]model.Voucher, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	whereClause := ""
	if activeOnly {
		whereClause = "WHERE is_active = 1"
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM vouchers " + whereClause
	err := s.DB.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get vouchers
	query := `
		SELECT id, code, description, type, value, min_order, max_discount,
			usage_limit, used_count, start_date, end_date, is_active, created_at
		FROM vouchers
		` + whereClause + `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var vouchers []model.Voucher
	for rows.Next() {
		var v model.Voucher
		err := rows.Scan(
			&v.ID, &v.Code, &v.Description, &v.Type, &v.Value,
			&v.MinOrder, &v.MaxDiscount, &v.UsageLimit, &v.UsedCount,
			&v.StartDate, &v.EndDate, &v.IsActive, &v.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		vouchers = append(vouchers, v)
	}

	return vouchers, total, nil
}

// UpdateVoucher cập nhật voucher
func (s *VoucherService) UpdateVoucher(voucher *model.Voucher) error {
	query := `
		UPDATE vouchers
		SET description = ?, type = ?, value = ?, min_order = ?, max_discount = ?,
			usage_limit = ?, start_date = ?, end_date = ?, is_active = ?
		WHERE id = ?
	`

	_, err := s.DB.Exec(query,
		voucher.Description,
		voucher.Type,
		voucher.Value,
		voucher.MinOrder,
		voucher.MaxDiscount,
		voucher.UsageLimit,
		voucher.StartDate,
		voucher.EndDate,
		voucher.IsActive,
		voucher.ID,
	)

	return err
}

// DeleteVoucher xóa voucher (soft delete - set is_active = false)
func (s *VoucherService) DeleteVoucher(id int64) error {
	_, err := s.DB.Exec("UPDATE vouchers SET is_active = 0 WHERE id = ?", id)
	return err
}
