package conf

import (
	"github.com/spf13/viper"
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

func TestJson(t *testing.T) {
	var example = []byte(`
{"abc": "good"}
`)

	LoadConfigByte(example, "json")
	t.Log(viper.GetString("abc"))
	t.Log(GetEnv("SY_REDIS_CACHE_HOST_PROXY"))
}

func TestEnv(t *testing.T) {
	t.Log(GetEnv("SY_REDIS_CACHE_HOST_PROXY"))
}

