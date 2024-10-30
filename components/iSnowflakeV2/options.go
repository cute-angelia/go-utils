package iSnowflakeV2

import "github.com/yitter/idgenerator-go/idgen"

const (
	defaultMethod            = 1             // 雪花计算方法,（1-漂移算法|2-传统算法），默认1
	defaultWorkerId          = 1             // 机器码，必须由外部设定，最大值 2^WorkerIdBitLength-1
	defaultBaseTime          = 1730291699000 // 基础时间（ms单位），不能超过当前系统时间
	defaultWorkerIdBitLength = 6             // 机器码位长，默认值6，取值范围 [1, 15]（要求：序列数位长+机器码位长不超过22）
	defaultSeqBitLength      = 6             // 序列数位长，默认值6，取值范围 [3, 21]（要求：序列数位长+机器码位长不超过22）
	defaultMaxSeqNumber      = 0             // 最大序列数（含），设置范围 [MinSeqNumber, 2^SeqBitLength-1]，默认值0，表示最大序列数取最大值（2^SeqBitLength-1]）
	defaultMinSeqNumber      = 5             // 最小序列数（含），默认值5，取值范围 [5, MaxSeqNumber]，每毫秒的前5个序列数对应编号0-4是保留位，其中1-4是时间回拨相应预留位，0是手工新值预留位
	defaultTopOverCostCount  = 2000          // 最大漂移次数（含），默认2000，推荐范围500-10000（与计算能力有关）
)

type OptionFuc func(o *options)

type options struct {
	method            uint16
	workerId          uint16
	baseTime          int64
	workerIdBitLength byte
	seqBitLength      byte
	maxSeqNumber      uint32
	minSeqNumber      uint32
	topOverCostCount  uint32

	idGenOption idgen.IdGeneratorOptions
}

func defaultOptions() *options {
	return &options{
		method:            defaultMethod,
		workerId:          defaultWorkerId,
		baseTime:          defaultBaseTime,
		workerIdBitLength: defaultWorkerIdBitLength,
		seqBitLength:      defaultSeqBitLength,
		maxSeqNumber:      defaultMaxSeqNumber,
		minSeqNumber:      defaultMinSeqNumber,
		topOverCostCount:  defaultTopOverCostCount,
	}
}

func WithMethod(methodId uint16) OptionFuc {
	return func(o *options) { o.method = methodId }
}

func WithWorkerId(workerId uint16) OptionFuc {
	return func(o *options) { o.workerId = workerId }
}

func WithBaseTime(baseTime int64) OptionFuc {
	return func(o *options) {
		o.baseTime = baseTime
	}
}

func WithWorkerIdBitLength(WorkerIdBitLength byte) OptionFuc {
	return func(o *options) {
		o.workerIdBitLength = WorkerIdBitLength
	}
}

func WithSeqBitLength(SeqBitLength byte) OptionFuc {
	return func(o *options) {
		o.seqBitLength = SeqBitLength
	}
}

func WithMaxSeqNumber(MaxSeqNumber uint32) OptionFuc {
	return func(o *options) {
		o.maxSeqNumber = MaxSeqNumber
	}
}

func WithMinSeqNumber(MinSeqNumber uint32) OptionFuc {
	return func(o *options) {
		o.minSeqNumber = MinSeqNumber
	}
}

func WithTopOverCostCount(TopOverCostCount uint32) OptionFuc {
	return func(o *options) {
		o.topOverCostCount = TopOverCostCount
	}
}
