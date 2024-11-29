package request_dto

type GetAllRequest struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}
