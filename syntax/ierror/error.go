package ierror

import "fmt"

type Code int32 //错误码

type Error struct {
	Message string
	Code    Code
}

// New returns an error object for the code, sendMessage.
func New(code Code, message string) error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func NewCode(code Code) error {
	return &Error{
		Code:    code,
		Message: code.String(),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

func (e *Error) Is(tgt error) bool {
	target, ok := tgt.(*Error)
	if !ok {
		return false
	}
	if e == nil || tgt == nil {
		return e == tgt
	}

	return e.Code == target.Code
}