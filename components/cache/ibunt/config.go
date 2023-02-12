package ibunt

// config options
type config struct {
	Name   string
	DbFile string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		Name:   "cache",
		DbFile: "cache_bunt.db",
	}
}
