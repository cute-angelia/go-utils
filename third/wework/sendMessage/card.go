package sendMessage

import "encoding/json"

type CardMessage struct {
	Touser  string `json:"touser"`
	Toparty string `json:"toparty"`
	Totag   string `json:"totag"`
	Msgtype string `json:"msgtype"`
	Agentid int    `json:"agentid"`
	News    struct {
		Articles []CardMessageArticle `json:"articles"`
	} `json:"news"`
	EnableIDTrans          int `json:"enable_id_trans"`
	EnableDuplicateCheck   int `json:"enable_duplicate_check"`
	DuplicateCheckInterval int `json:"duplicate_check_interval"`
}

type CardMessageArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Picurl      string `json:"picurl"`
	Appid       string `json:"appid"`
	Pagepath    string `json:"pagepath"`
}

func (receiver CardMessage) GetMessage() []byte {
	v, _ := json.Marshal(receiver)
	return v
}
