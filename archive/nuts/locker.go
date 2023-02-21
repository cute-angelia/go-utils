/*
* 限制类
  次数限制
  每日次数限制
*/
package nuts

type LockerOpts struct {
	Limit int   // 限制次数
	Today bool  // 是否限制每日
	Uid   int32 // 针对某人
}

type LockerOpt func(opts *LockerOpts)

func NewLockerOpt(opts ...LockerOpt) LockerOpts {
	var sopt LockerOpts
	for _, opt := range opts {
		opt(&sopt)
	}
	if sopt.Limit == 0 {
		sopt.Limit = 1
	}
	return sopt
}

func WithLimit(limit int) LockerOpt {
	return func(options *LockerOpts) {
		options.Limit = limit
	}
}

func WithToday(today bool) LockerOpt {
	return func(options *LockerOpts) {
		options.Today = today
	}
}

func WithUid(uid int32) LockerOpt {
	return func(options *LockerOpts) {
		options.Uid = uid
	}
}
