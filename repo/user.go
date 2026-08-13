package repo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/shohann/golang-ecommerce-api/apperr"
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
	user.Role = "customer"

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
		return nil, apperr.Internal("create user", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&userId); err != nil {
			return nil, apperr.Internal("scan user id", err)
		}
	}

	user.ID = userId

	return &user, nil
}

func (r *userRepo) FindUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := `
	  SELECT id, full_name, email, password_hash, role
      FROM users
      WHERE email = $1
      LIMIT 1
	`

	err := r.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperr.NotFound("user not found")
		}
		return nil, apperr.Internal("find auth user", err)
	}

	return &user, nil
}

func (r *userRepo) FindUserById(id int) (*domain.User, error) {
	var user domain.User

	query := `
	  SELECT id, full_name, email, password_hash, role
      FROM users
      WHERE id = $1
      LIMIT 1
	`

	err := r.db.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperr.NotFound("user not found")
		}
		return nil, apperr.Internal("find auth user", err)
	}

	return &user, nil

}

func (r *userRepo) CheckUniqueUser(email string) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM users 
			WHERE email = $1
		);
	`
	err := r.db.Get(&exists, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, apperr.Internal("check unique user", err)
	}

	return exists, nil
}
