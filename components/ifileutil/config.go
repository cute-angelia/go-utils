package ifileutil

const ComponentName = "Component.IfileUtil"

type config struct {
	Dir        string
	DirInclude []string
	DirDeclude []string
	ExtInclude []string
	ExtDeclude []string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{}
}
