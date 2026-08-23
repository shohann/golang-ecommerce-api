package repo

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/shohann/golang-ecommerce-api/apperr"
	"github.com/shohann/golang-ecommerce-api/category"
	"github.com/shohann/golang-ecommerce-api/domain"
)

type CategoryRepo interface {
	category.CategoryRepo
}

type categoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) CategoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) Create(cat domain.Category) (*domain.Category, error) {
	query := `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING id, name
	`

	var created domain.Category
	err := r.db.QueryRowx(query, cat.Name).StructScan(&created)
	if err != nil {
		return nil, mapCategoryDBError("create category", err)
	}

	return &created, nil
}

func (r *categoryRepo) FindByID(id int64) (*domain.Category, error) {
	var cat domain.Category
	query := `
		SELECT id, name
		FROM categories
		WHERE id = $1
		LIMIT 1
	`

	err := r.db.Get(&cat, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, apperr.Internal("find category", err)
	}

	return &cat, nil
}

func (r *categoryRepo) Update(cat domain.Category) (*domain.Category, error) {
	query := `
		UPDATE categories
		SET name = $1
		WHERE id = $2
		RETURNING id, name
	`

	var updated domain.Category
	err := r.db.QueryRowx(query, cat.Name, cat.ID).StructScan(&updated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, mapCategoryDBError("update category", err)
	}

	return &updated, nil
}

func (r *categoryRepo) Delete(id int64) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return mapCategoryDBError("delete category", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return apperr.Internal("delete category rows", err)
	}
	if rows == 0 {
		return apperr.NotFound("category not found")
	}

	return nil
}

func (r *categoryRepo) List(limit, offset int64) ([]domain.Category, error) {
	var categories []domain.Category

	query := `
		SELECT id, name
		FROM categories
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	err := r.db.Select(&categories, query, limit, offset)
	if err != nil {
		return nil, apperr.Internal("list categories", err)
	}

	if categories == nil {
		categories = []domain.Category{}
	}

	return categories, nil
}

func (r *categoryRepo) Count() (int64, error) {
	var count int64
	err := r.db.Get(&count, `SELECT COUNT(*) FROM categories`)
	if err != nil {
		return 0, apperr.Internal("count categories", err)
	}
	return count, nil
}

func (r *categoryRepo) ExistsByName(name string, excludeID int64) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM categories
			WHERE name = $1 AND id <> $2
		)
	`

	err := r.db.Get(&exists, query, name, excludeID)
	if err != nil {
		return false, apperr.Internal("check category name", err)
	}

	return exists, nil
}

func mapCategoryDBError(msg string, err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return apperr.Conflict("category name already exists")
		case "23503":
			return apperr.Conflict("category is in use")
		}
	}
	return apperr.Internal(msg, err)
}
