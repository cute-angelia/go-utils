package iredis

// config options
type config struct {
	Name     string
	Server   string
	Password string
	DbIndex  int
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		Name:     "default",
		Server:   "127.0.0.1:6379",
		Password: "",
		DbIndex:  0,
	}
}
