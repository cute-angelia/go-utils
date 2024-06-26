package iupload

import "time"

// config options
type config struct {

	// 上传文件夹
	UploadDirectory string
	// 上传图片限制
	UploadImageSize int64
	// 上传视频限制
	UploadVideoSize int64
	// 上传图片扩展
	UploadImageExt []string
	// 上传视频扩展
	UploadVideoExt []string

	Debug   bool          //  打印日志
	Timeout time.Duration // 超时时间

	ReplaceMode int // 替换模式， 1跳过， 2覆盖  3保留两者
}

const (
	ReplaceModeIgnore  = 1
	ReplaceModeReplace = 2
	ReplaceModeTwo     = 3
)

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		Debug:           false,
		ReplaceMode:     2,
		UploadImageSize: 1024 * 1024 * 10,
		UploadVideoSize: 1024 * 1024 * 1024 * 10,
		UploadImageExt:  []string{"png", "jpg", "jpeg", "gif", "ico", "bmp"},
		UploadVideoExt:  []string{"mp4", "mp3", "avi", "flv", "rmvb", "mov"},
	}
}
