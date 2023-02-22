package routercache

import (
	"github.com/cute-angelia/go-utils/components/caches"
	"net/http"
	"time"
)

const PackageName = "component.utils.http.routercache"

// config options
type config struct {
	Store              caches.Cache
	Ttl                time.Duration
	RefreshKey         string
	Methods            []string
	StatusCodeFilter   func(int) bool
	WriteExpiresHeader bool
	CustomKey          string
	PrintLog           bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		Ttl:                time.Minute * 10,
		RefreshKey:         "xfresh",
		Methods:            []string{http.MethodPost, http.MethodGet},
		WriteExpiresHeader: true,
		CustomKey:          "",
		PrintLog:           true,
	}
}
