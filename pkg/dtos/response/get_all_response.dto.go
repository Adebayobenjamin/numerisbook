package response_dto

type GetAllResponse[T any] struct {
	Data        []T  `json:"data"`
	TotalCount  int  `json:"total_count"`
	TotalPages  int  `json:"total_pages"`
	CurrentPage int  `json:"current_page"`
	Limit       int  `json:"limit"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}
