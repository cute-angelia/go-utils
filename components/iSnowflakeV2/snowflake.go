package iSnowflakeV2

import "github.com/yitter/idgenerator-go/idgen"

type Snowflake struct {
	opts *options
	// sfg  singleflight.Group // "golang.org/x/sync/singleflight"
}

func NewSnowflake(opts ...OptionFuc) *Snowflake {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	o.idGenOption = idgen.IdGeneratorOptions{
		Method:            o.method,
		WorkerId:          o.workerId,
		BaseTime:          o.baseTime,
		WorkerIdBitLength: o.workerIdBitLength,
		SeqBitLength:      o.seqBitLength,
		MaxSeqNumber:      o.maxSeqNumber,
		MinSeqNumber:      o.minSeqNumber,
		TopOverCostCount:  o.topOverCostCount,
	}
	idgen.SetIdGenerator(&o.idGenOption)

	return &Snowflake{
		opts: o,
	}
}

func (that *Snowflake) GetId() int64 {
	return idgen.NextId()
}
