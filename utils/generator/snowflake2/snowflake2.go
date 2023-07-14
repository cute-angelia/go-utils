package snowflake2

import (
	"github.com/yitter/idgenerator-go/idgen"
	"time"
)

func InitSnowflake() {
	// 创建 IdGeneratorOptions 对象，可在构造函数中输入 WorkerId：
	var options = idgen.NewIdGeneratorOptions(1)
	// options.WorkerIdBitLength = 10  // 默认值6，限定 WorkerId 最大值为2^6-1，即默认最多支持64个节点。
	options.SeqBitLength = 6                                                   //           // 默认值6，限制每毫秒生成的ID个数。若生成速度超过5万个/秒，建议加大 SeqBitLength 到 10。                                     // 默认值6，限制每毫秒生成的ID个数。若生成速度超过5万个/秒，建议加大 SeqBitLength 到 10。
	options.BaseTime = time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC).UnixMilli() // 如果要兼容老系统的雪花算法，此处应设置为老系统的BaseTime。
	// ...... 其它参数参考 IdGeneratorOptions 定义。
	idgen.SetIdGenerator(options)
}

func SetOption(workerId uint16, seqBitLength byte) {
	var options = idgen.NewIdGeneratorOptions(workerId)
	// options.WorkerIdBitLength = 10  // 默认值6，限定 WorkerId 最大值为2^6-1，即默认最多支持64个节点。
	options.SeqBitLength = seqBitLength
	options.BaseTime = time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC).UnixMilli() // 如果要兼容老系统的雪花算法，此处应设置为老系统的BaseTime。
	idgen.SetIdGenerator(options)
}

func GetNextId() int64 {
	return idgen.NextId()
}
