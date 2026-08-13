package user

import "github.com/shohann/golang-ecommerce-api/domain"

type UserResponse struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	ID          int    `json:"id"`
	AccessToken string `json:"access_token"`
}

func ToUserResponse(u *domain.User) UserResponse {
	if u == nil {
		return UserResponse{}
	}

	return UserResponse{
		ID:       u.ID,
		FullName: u.FullName,
		Email:    u.Email,
		Role:     u.Role,
	}
}

func ToUserResponses(users []domain.User) []UserResponse {
	result := make([]UserResponse, 0, len(users))
	for i := range users {
		result = append(result, ToUserResponse(&users[i]))
	}
	return result
}
