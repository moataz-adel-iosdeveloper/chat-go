package models

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func SuccessResponse(message string, data any) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string, err any) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
}
