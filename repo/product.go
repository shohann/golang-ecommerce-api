package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shohann/golang-ecommerce-api/apperr"
	"github.com/shohann/golang-ecommerce-api/domain"
	"github.com/shohann/golang-ecommerce-api/product"
)

type ProductRepo interface {
	product.ProductRepo
}

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) ProductRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Create(product domain.Product) (*domain.Product, error) {
	query := `
        INSERT INTO products (
            category_id,
            name,
            description,
            price,
            stock,
            image_url,
            is_active
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING
            id, category_id, name, description, price, stock,
            image_url, is_active, created_at, updated_at
    `

	err := r.db.QueryRow(
		query,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ImageURL,
		product.IsActive,
	).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, apperr.Internal("create_product", err)
	}

	return &product, nil
}

func (r *productRepo) List(limit, offset int64) ([]domain.Product, error) {
	var products []domain.Product

	query := `
	SELECT *
	FROM products
	ORDER BY id DESC
	LIMIT $1 OFFSET $2
	`

	err := r.db.Select(&products, query, limit, offset)

	fmt.Println(err)

	if err != nil {
		return nil, apperr.Internal("list_product", err)
	}

	if products == nil {
		products = []domain.Product{}
	}

	return products, nil
}

func (r *productRepo) Count() (int64, error) {
	var count int64

	query := `SELECT COUNT(*) FROM products`

	err := r.db.Get(&count, query)
	if err != nil {
		return 0, apperr.Internal("count_products", err)
	}

	fmt.Println(err)

	return count, nil
}

func (r *productRepo) Delete(id int64) error {
	fmt.Println(id)
	query := `DELETE FROM products WHERE id = $1`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return apperr.Internal("delete_product", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return apperr.Internal("delete_product_rows_affected", err)
	}

	if rowsAffected == 0 {
		return apperr.NotFound("product_not_found")
	}

	return nil
}
