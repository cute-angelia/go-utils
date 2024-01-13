package ijson

import (
	"encoding/json"
	"log"
	"testing"
)

func TestJson(t *testing.T) {
	var str = `{
"ToUserName": "gh_c686a729890a",
"FromUserName": "o04qa5RCXTE48QzCg5tJbLB-jpx0",
"CreateTime": 1641265084,
"MsgType": "text",
"Content": "1",
"MsgId": 23497265183919412
}`
	var params map[string]string
	err := Unmarshal([]byte(str), &params)
	log.Println(err)
	log.Println(Pretty(params))

	var params2 map[string]string
	err2 := json.Unmarshal([]byte(str), &params2)
	log.Println(err2)
	log.Println(Pretty(params2))

	var numpool []int
	log.Println(UnmarshalSlice([]byte("[2,3,4]"), &numpool))
	log.Println(numpool)

	var numpool3 [][]int
	log.Println(UnmarshalSlice([]byte("[[2500,2500,2000,2000,1000],[2500,2500,2000,2000,1000],[7500,2000,500],[7500,2000,500],[7500,2000,500],[9500,500]]"), &numpool3))
	log.Println(numpool3)

	var numpool2 []string
	log.Println(UnmarshalSlice([]byte(`["gold","exp","coinMinting1","coinMinting2","soulStone1","soulStone2"]`), &numpool2))
	log.Println(numpool2)
}
