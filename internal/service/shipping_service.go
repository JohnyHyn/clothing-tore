package service

import (
	"clothing-store/internal/model"
	"database/sql"
	"errors"
)

type ShippingService struct {
	DB *sql.DB
}

// CreateShipping tạo thông tin vận chuyển cho order
func (s *ShippingService) CreateShipping(shipping *model.Shipping) error {
	query := `
		INSERT INTO shippings (order_id, method, provider, tracking_code, fee, status,
			address, city, district, ward, receiver_name, receiver_phone, note, estimated_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.DB.Exec(query,
		shipping.OrderID,
		shipping.Method,
		shipping.Provider,
		shipping.TrackingCode,
		shipping.Fee,
		shipping.Status,
		shipping.Address,
		shipping.City,
		shipping.District,
		shipping.Ward,
		shipping.ReceiverName,
		shipping.ReceiverPhone,
		shipping.Note,
		shipping.EstimatedDate,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	shipping.ID = id
	return nil
}

// GetShippingByOrderID lấy thông tin shipping theo order ID
func (s *ShippingService) GetShippingByOrderID(orderID int64) (*model.Shipping, error) {
	shipping := &model.Shipping{}

	query := `
		SELECT id, order_id, method, provider, tracking_code, fee, status,
			address, city, district, ward, receiver_name, receiver_phone, note,
			estimated_date, delivered_at, created_at, updated_at
		FROM shippings
		WHERE order_id = ?
	`

	var deliveredAt, estimatedDate sql.NullTime
	var ward, note sql.NullString

	err := s.DB.QueryRow(query, orderID).Scan(
		&shipping.ID,
		&shipping.OrderID,
		&shipping.Method,
		&shipping.Provider,
		&shipping.TrackingCode,
		&shipping.Fee,
		&shipping.Status,
		&shipping.Address,
		&shipping.City,
		&shipping.District,
		&ward,
		&shipping.ReceiverName,
		&shipping.ReceiverPhone,
		&note,
		&estimatedDate,
		&deliveredAt,
		&shipping.CreatedAt,
		&shipping.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if ward.Valid {
		shipping.Ward = ward.String
	}
	if note.Valid {
		shipping.Note = note.String
	}
	if estimatedDate.Valid {
		shipping.EstimatedDate = estimatedDate.Time
	}
	if deliveredAt.Valid {
		shipping.DeliveredAt = deliveredAt.Time
	}

	return shipping, nil
}

// UpdateShippingStatus cập nhật trạng thái vận chuyển
func (s *ShippingService) UpdateShippingStatus(shippingID int64, status, location, note string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update shipping status
	query := "UPDATE shippings SET status = ?, updated_at = NOW() WHERE id = ?"

	if status == "delivered" {
		query = "UPDATE shippings SET status = ?, delivered_at = NOW(), updated_at = NOW() WHERE id = ?"
	}

	_, err = tx.Exec(query, status, shippingID)
	if err != nil {
		return err
	}

	// Add to history
	_, err = tx.Exec(`
		INSERT INTO shipping_history (shipping_id, status, location, note)
		VALUES (?, ?, ?, ?)
	`, shippingID, status, location, note)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetShippingHistory lấy lịch sử vận chuyển
func (s *ShippingService) GetShippingHistory(shippingID int64) ([]model.ShippingHistory, error) {
	query := `
		SELECT id, shipping_id, status, location, note, created_at
		FROM shipping_history
		WHERE shipping_id = ?
		ORDER BY created_at DESC
	`

	rows, err := s.DB.Query(query, shippingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []model.ShippingHistory
	for rows.Next() {
		var h model.ShippingHistory
		var location, note sql.NullString

		err := rows.Scan(&h.ID, &h.ShippingID, &h.Status, &location, &note, &h.CreatedAt)
		if err != nil {
			return nil, err
		}

		if location.Valid {
			h.Location = location.String
		}
		if note.Valid {
			h.Note = note.String
		}

		history = append(history, h)
	}

	return history, nil
}

// TrackShipping theo dõi đơn hàng theo tracking code
func (s *ShippingService) TrackShipping(trackingCode string) (*model.Shipping, []model.ShippingHistory, error) {
	shipping := &model.Shipping{}

	query := `
		SELECT id, order_id, method, provider, tracking_code, fee, status,
			address, city, district, ward, receiver_name, receiver_phone, note,
			estimated_date, delivered_at, created_at, updated_at
		FROM shippings
		WHERE tracking_code = ?
	`

	var deliveredAt, estimatedDate sql.NullTime
	var ward, note sql.NullString

	err := s.DB.QueryRow(query, trackingCode).Scan(
		&shipping.ID,
		&shipping.OrderID,
		&shipping.Method,
		&shipping.Provider,
		&shipping.TrackingCode,
		&shipping.Fee,
		&shipping.Status,
		&shipping.Address,
		&shipping.City,
		&shipping.District,
		&ward,
		&shipping.ReceiverName,
		&shipping.ReceiverPhone,
		&note,
		&estimatedDate,
		&deliveredAt,
		&shipping.CreatedAt,
		&shipping.UpdatedAt,
	)

	if err != nil {
		return nil, nil, err
	}

	if ward.Valid {
		shipping.Ward = ward.String
	}
	if note.Valid {
		shipping.Note = note.String
	}
	if estimatedDate.Valid {
		shipping.EstimatedDate = estimatedDate.Time
	}
	if deliveredAt.Valid {
		shipping.DeliveredAt = deliveredAt.Time
	}

	// Get history
	history, err := s.GetShippingHistory(shipping.ID)
	if err != nil {
		return shipping, nil, err
	}

	return shipping, history, nil
}

// CalculateShippingFee tính phí ship dựa trên method và địa chỉ
func (s *ShippingService) CalculateShippingFee(method, city string) (float64, error) {
	// Đây là logic đơn giản, trong thực tế có thể gọi API của đơn vị vận chuyển
	baseFee := map[string]float64{
		"standard": 30000,
		"express":  50000,
		"same_day": 80000,
	}

	fee, ok := baseFee[method]
	if !ok {
		return 0, errors.New("invalid shipping method")
	}

	// Điều chỉnh phí theo thành phố (ví dụ)
	remoteCities := []string{"Hà Giang", "Cao Bằng", "Lai Châu", "Lào Cai"}
	for _, remoteCity := range remoteCities {
		if city == remoteCity {
			fee += 20000 // Thêm phí cho tỉnh xa
			break
		}
	}

	return fee, nil
}

// UpdateTrackingCode cập nhật mã vận đơn
func (s *ShippingService) UpdateTrackingCode(shippingID int64, trackingCode string) error {
	_, err := s.DB.Exec(
		"UPDATE shippings SET tracking_code = ?, updated_at = NOW() WHERE id = ?",
		trackingCode, shippingID,
	)
	return err
}
