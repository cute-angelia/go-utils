package sendMessage

import "encoding/json"

type TextMsg struct {
	Touser  string `json:"touser,omitempty"`
	Toparty string `json:"toparty,omitempty"`
	Totag   string `json:"totag,omitempty"`
	Msgtype string `json:"msgtype"`
	Agentid int    `json:"agentid"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Safe                   int `json:"safe,omitempty"`
	EnableIDTrans          int `json:"enable_id_trans,omitempty"`
	EnableDuplicateCheck   int `json:"enable_duplicate_check,omitempty"`
	DuplicateCheckInterval int `json:"duplicate_check_interval,omitempty"`
}

func (receiver TextMsg) GetMessage() []byte {
	v, _ := json.Marshal(receiver)
	return v
}
