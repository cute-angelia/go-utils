package wework

import (
	"github.com/cute-angelia/go-utils/v2/components/conf"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	conf.LoadConfigFile("./config.toml")
	Load("wework")
	iapp, _ := GetApp("1000012")
	accessToken, err := iapp.AccessToken.GetAccessTokenHandler(iapp)
	t.Log(accessToken, err)
}
