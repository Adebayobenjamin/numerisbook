package common

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func BuildSuccessResponse(message string, data interface{}) Response {
	return Response{Status: "success", Message: message, Data: data}
}

func BuildErrorResponse(err interface{}, message string) Response {
	return Response{Status: "error", Message: message, Error: err}

}
