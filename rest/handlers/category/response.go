package category

import "github.com/shohann/golang-ecommerce-api/domain"

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func ToCategoryResponse(c *domain.Category) CategoryResponse {
	if c == nil {
		return CategoryResponse{}
	}

	return CategoryResponse{
		ID:   c.ID,
		Name: c.Name,
	}
}

func ToCategoryResponses(categories []domain.Category) []CategoryResponse {
	result := make([]CategoryResponse, 0, len(categories))
	for i := range categories {
		result = append(result, ToCategoryResponse(&categories[i]))
	}
	return result
}
