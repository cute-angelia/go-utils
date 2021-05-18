package download_img

// config options
type config struct {
	Width  int    // 最小宽度 
	Height int    // 最小高度
	Rename bool   // 是否重命名
	Dest   string // 下载路径

	Cookie     string
	UserAgent  string
	DefaultExt string // 默认扩展
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		Dest:       "./",
		Width:      0,
		Height:     0,
		Rename:     false,
		Cookie:     "",
		UserAgent:  "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
		DefaultExt: "",
	}
}
