package isnowflake

import "github.com/yitter/idgenerator-go/idgen"

type Component struct {
	config *config
}

func newComponent(config *config) *Component {
	comp := &Component{}
	comp.config = config

	options := idgen.IdGeneratorOptions{
		Method:            config.Method,
		WorkerId:          config.Method,
		BaseTime:          config.BaseTime,
		WorkerIdBitLength: config.WorkerIdBitLength,
		SeqBitLength:      config.SeqBitLength,
		MaxSeqNumber:      config.MaxSeqNumber,
		MinSeqNumber:      config.MinSeqNumber,
		TopOverCostCount:  config.TopOverCostCount,
	}
	idgen.SetIdGenerator(&options)

	return comp
}

func (that *Component) GetId() int64 {
	return idgen.NextId()
}
