package uid

type IdCryptOptions struct {
	Dict64 [64]uint64
	Key64  uint64
	Dict32 [32]uint32
	Key32  uint32
}

type IdCryptOption func(options *IdCryptOptions)

func NewIdCryptOptions(opts ...IdCryptOption) IdCryptOptions {
	var initOpts IdCryptOptions

	//if &initOpts.Dict64 == nil {
	initOpts.Dict64 = [64]uint64{
		0, 1, 63, 3, 4, 55, 6, 7, 8, 9, 60,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
		31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
		41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
		51, 52, 53, 54, 5, 56, 57, 58, 59, 10,
		61, 62, 2,
	}
	//}

	initOpts.Key64 = uint64(2018101321)
	initOpts.Dict32 = [32]uint32{
		0, 1, 22, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 2, 23, 24, 25, 26, 27, 28, 29, 30,
		31,
	}
	initOpts.Key32 = uint32(1314520)

	for _, opt := range opts {
		opt(&initOpts)
	}

	return initOpts
}

func WithDict64(Dict [64]uint64) IdCryptOption {
	return func(options *IdCryptOptions) {
		options.Dict64 = Dict
	}
}

func WithKey64(key uint64) IdCryptOption {
	return func(options *IdCryptOptions) {
		options.Key64 = key
	}
}

func WithDict32(dict [32]uint32) IdCryptOption {
	return func(options *IdCryptOptions) {
		options.Dict32 = dict
	}
}

func WithKey32(key uint32) IdCryptOption {
	return func(options *IdCryptOptions) {
		options.Key32 = key
	}
}
