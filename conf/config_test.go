package conf

import (
	"bytes"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToml(t *testing.T) {
	var example = []byte(`
[common]
debug=true
`)

	MustLoadConfigByte(example, "toml")
	t.Log(viper.GetBool("common.debug"))
}

func TestMerge(t *testing.T) {
	var example = []byte(`
[common]
debug=true
`)

	var example2 = []byte(`
[common]
debug=false
hello="world"
[key]
hello="keyworld"
`)

	MustLoadConfigByte(example, "toml")

	// 合并
	MergeConfig(bytes.NewReader(example2))

	t.Log(viper.GetBool("common.debug"))
	t.Log(viper.GetString("common.hello"))
	t.Log(viper.GetString("key.hello"))
}

func TestJson(t *testing.T) {
	var example = []byte(`
{"abc": "good"}
`)

	LoadConfigByte(example, "json")
	t.Log(viper.GetString("abc"))
	t.Log(GetEnv("SY_REDIS_CACHE_HOST_PROXY"))
}

func TestEnv(t *testing.T) {
	assert.Equal(t, "192.168.1.234", GetEnv("SY_REDIS_CACHE_HOST_PROXY"))
}
