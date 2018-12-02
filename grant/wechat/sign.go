package wechat

import (
	"sort"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"fmt"
)

const DEFAULT_WECHAT_KEY = "19dfg06250bzxc9247ec02edce69f6a2d"

type Sign struct {
	apiKey string
}

// 创建
func NewSign() *Sign {
	return &Sign{
		apiKey: DEFAULT_WECHAT_KEY,
	}
}

func (c *Sign) SetApiKey(apiKey string) {
	c.apiKey = apiKey
}

// 验证签名
func (c *Sign) ValidSign(signIn string, signOut string) bool {
	return signIn == signOut
}

/**
	Signature 签名算法
	query = r.URL.Query()
 */
func (c *Sign) Signature(query map[string][]string) string {
	// 排序
	keys := make([]string, 0, len(query))
	for k := range query {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接
	querypad := []string{}
	for _, k := range keys {
		v := strings.Join(query[k], ",")
		querypad = append(querypad, fmt.Sprintf("%s=%s", k, v))
	}
	querypad = append(querypad, "key="+c.apiKey)
	stringA := strings.Join(querypad, "&")

	//log.Println("sign str:", stringA)
	// 计算
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(stringA))
	sign := hex.EncodeToString(md5Ctx.Sum(nil))
	return strings.ToUpper(sign)
}
