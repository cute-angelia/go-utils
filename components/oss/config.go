package oss

// config options
type config struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string

	BucketName string
	BucketHost string

	ProxySocks5 string

	Debug bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		Debug: false,
	}
}
