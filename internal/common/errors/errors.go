package errors

type ApiError struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Code string `json:"code,omitempty"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func NewApiError(status int, message string) *ApiError {
	return &ApiError{
		Status: status,
		Message: message,
	}
}