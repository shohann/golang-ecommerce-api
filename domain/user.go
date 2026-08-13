package domain

type User struct {
	ID       int    `json:"id" db:"id"`
	FullName string `json:"full_name" db:"full_name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password_hash"`
	Role     string `json:"role" db:"role"`
}

type LoginResult struct {
	ID          int    `json:"id"`
	AccessToken string `json:"access_token"`
}
