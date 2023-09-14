package media

import (
	"net/url"
	"strings"
)

// GetLargeImg 获取大图
func GetLargeImg(imgURL string) string {
	if imgURL != "" {
		// Weibo
		if strings.Contains(imgURL, "sinaimg.cn") {
			u, _ := url.Parse(imgURL)
			parts := strings.Split(u.Path, "/")
			parts[1] = "large"
			imgURL = u.Scheme + "://" + u.Host + strings.Join(parts, "/")
		}

		// Twitter
		if strings.Contains(imgURL, "twimg.com") {
			if strings.Contains(imgURL, "?") {
				imgURL = imgURL[:strings.Index(imgURL, "?")]
			}
			imgURL += "?format=jpg&name=orig"
		}
	}

	return imgURL
}
