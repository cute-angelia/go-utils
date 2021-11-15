package ierror

import (
	"github.com/cute-angelia/go-utils/syntax/ijson"
)

type Err struct {
	Code int
	Msg  string
	Data interface{}
}

func (e *Err) Error() string {
	err, _ := ijson.Marshal(e)
	return string(err)
}
func New(code int, msg string, data interface{}) *Err {
	return &Err{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
