package user

type ReqCreateUser struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReqLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
