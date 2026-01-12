package service

import (
	"clothing-store/internal/model"
	"database/sql"
)

type ProductService struct {
	DB *sql.DB
}

func (s *ProductService) CreateProduct(p *model.Product) error {
	query := `
		INSERT INTO products
		(name, description, price, stock, category, image_url, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.DB.Exec(
		query,
		p.Name,
		p.Description,
		p.Price,
		p.Stock,
		p.Category,
		p.ImageURL,
		"active",
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = id
	return nil
}
func (s *ProductService) GetProducts(page, limit int) ([]model.Product, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, name, price, stock, category, image_url, status
		FROM products
		WHERE status = 'active'
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []model.Product{}

	for rows.Next() {
		var p model.Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.Category,
			&p.ImageURL,
			&p.Status,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
func (s *ProductService) GetProductByID(id int64) (*model.Product, error) {
	query := `
		SELECT id, name, description, price, stock, category, image_url, status
		FROM products
		WHERE id = ?
	`

	var p model.Product
	err := s.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.Category,
		&p.ImageURL,
		&p.Status,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
func (s *ProductService) UpdateProduct(id int64, p *model.Product) error {
	query := `
		UPDATE products
		SET name = ?, description = ?, price = ?, stock = ?, category = ?, image_url = ?, status = ?
		WHERE id = ?
	`

	_, err := s.DB.Exec(
		query,
		p.Name,
		p.Description,
		p.Price,
		p.Stock,
		p.Category,
		p.ImageURL,
		p.Status,
		id,
	)
	return err
}
func (s *ProductService) DeleteProduct(id int64) error {
	query := `
		UPDATE products
		SET status = 'inactive'
		WHERE id = ?
	`
	_, err := s.DB.Exec(query, id)
	return err
}
