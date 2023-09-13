package apiV2

type ApiError struct {
	Code    int
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}

func NewApiError(code int, message string) *ApiError {
	return &ApiError{
		Code:    code,
		Message: message,
	}
}

func ApiErrorMsg(err error) *ApiError {
	if err != nil {
		if e, ok := err.(*ApiError); ok {
			return e
		}
	}
	return &ApiError{
		Code:    -1,
		Message: err.Error(),
	}
}
