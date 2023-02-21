package uid

func NewIdCrypt64(opt ...IdCryptOption) IdCrypt64 {
	opts := NewIdCryptOptions(opt...)
	crypt := IdCrypt64{
		opts: opts,
	}
	return crypt
}

type IdCrypt64 struct {
	opts IdCryptOptions
}

func (self IdCrypt64) Encrypt(id uint64) uint64 {
	iid := uint64(id | (1 << 62))
	newId := uint64(0)
	for k, v := range self.opts.Dict64 {
		a := (iid >> uint64(63-v)) & 1
		newId |= a << uint64(63-k)
	}
	return uint64(newId ^ self.opts.Key64)
}

func (self IdCrypt64) Decrypt(id uint64) uint64 {
	iid := id ^ self.opts.Key64
	newId := uint64(0)
	for k, v := range self.opts.Dict64 {
		a := uint64(iid>>uint64(63-k)) & 1
		newId |= a << uint64(63-v)
	}
	return uint64(newId & (^uint64(1 << 62)))
}
