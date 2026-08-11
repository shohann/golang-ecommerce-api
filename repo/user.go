package repo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/shohann/golang-ecommerce-api/domain"
	"github.com/shohann/golang-ecommerce-api/user"
)

type UserRepo interface {
	user.UserRepo
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(user domain.User) (*domain.User, error) {
	query := `
	    INSERT INTO users (
            full_name, 
            email, 
            password_hash
        )
        VALUES (
            :full_name,
            :email, 
            :password_hash
        )
        RETURNING id
	`

	var userId int
	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		return nil, err
	}
	// defer rows.Close() // Always close sql.Rows to avoid connection leaks

	if rows.Next() {
		rows.Scan(&userId)
	}

	user.ID = userId

	return &user, nil
}

func (r *userRepo) FindAuthUser(email, pass string) (*domain.User, error) {
	var user domain.User
	query := `
	  SELECT id, full_name, email, password_hash, role
      FROM users
      WHERE email = $1 AND password_hash = $2
      LIMIT 1
	`

	err := r.db.Get(&user, query, email, pass)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // no matching user
		}

		return nil, err
	}

	return &user, nil
}
