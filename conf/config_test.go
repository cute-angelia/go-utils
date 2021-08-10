package conf

import (
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
