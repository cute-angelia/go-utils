package wework

import (
	"github.com/cute-angelia/go-utils/components/iredis"
	"github.com/cute-angelia/go-utils/utils/conf"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	conf.LoadConfigFile("./config.toml")
	iredis.Load("redis").Init()

	Load("wework").InitWeWork(iredis.GetRedisClient("cache"))

	iapp, _ := GetApp("1000012")
	accessToken, err := iapp.AccessToken.GetAccessTokenHandler(iapp)
	t.Log(accessToken, err)
}
