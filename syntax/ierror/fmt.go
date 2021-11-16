package ierror

import (
	"github.com/cute-angelia/go-utils/syntax/ijson"
)

type Err struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
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

func Decode(str string) Err {
	e := Err{}
	if err := ijson.Unmarshal([]byte(str), &e); err != nil {
		e.Code = -1
		e.Msg = str + " + 解析失败"
	}
	return e
}
