package apiV3

type ApiError struct {
	Code    int32
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}

func NewApiError(code int32, message string) *ApiError {
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
