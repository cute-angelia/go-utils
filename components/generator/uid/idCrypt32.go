package uid

func NewIdCrypt32(opt ...IdCryptOption) IdCrypt32 {
	opts := NewIdCryptOptions(opt...)
	crypt := IdCrypt32{
		opts: opts,
	}
	return crypt
}

type IdCrypt32 struct {
	opts IdCryptOptions
}

func (self IdCrypt32) Encrypt(id uint32) uint32 {
	iid := uint32(id | (1 << 30))
	newId := uint32(0)
	for k, v := range self.opts.Dict32 {
		a := (iid >> uint32(31-v)) & 1
		newId |= a << uint32(31-k)
	}
	return uint32(newId ^ self.opts.Key32)
}

func (self IdCrypt32) Decrypt(id uint32) uint32 {
	iid := id ^ self.opts.Key32
	newId := uint32(0)
	for k, v := range self.opts.Dict32 {
		a := uint32(iid>>uint32(31-k)) & 1
		newId |= a << uint32(31-v)
	}
	return uint32(newId & (^uint32(1 << 30)))
}
