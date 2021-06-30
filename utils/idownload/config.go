package idownload

// config options
type config struct {
	Width  int // 图片-最小宽度
	Height int // 图片-最小高度

	Rename     bool   // 文件-是否重命名
	Dest       string // 文件-下载路径
	NamePrefix string // 文件-命名前缀
	DefaultExt string // 文件-后缀属性，有些文件没有后缀

	ProxyHttp   string // 代理 http://ip:port
	ProxySocks5 string // 代理 ip:port
	Cookie      string // cookie
	UserAgent   string // user-agent
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		Dest:        "./",
		Width:       0,
		Height:      0,
		Rename:      false,
		NamePrefix:  "",
		ProxyHttp:   "",
		ProxySocks5: "",
		Cookie:      "",
		UserAgent:   "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
		DefaultExt:  "",
	}
}
