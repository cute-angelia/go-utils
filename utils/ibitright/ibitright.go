package ibitright

import (
	"errors"
	"sync"
)

// 位管理
// 用于授权，或者一次性记录

const (
	MaxPosition = 32
)

var (
	ErrInvalidPosition = errors.New("invalid position")
)

// Position 位置 1 开始
type Position uint32

type ibitright struct {
	mu sync.Mutex
}

func NewBitRight() *ibitright {
	return &ibitright{}
}

// position Position // 设置第几位
// oldNum   uint32   // 存储的值，默认0

func (that *ibitright) Check(position Position, oldNum uint32) (bool, error) {
	if position < 1 || position > MaxPosition {
		return false, ErrInvalidPosition
	}

	that.mu.Lock()
	defer that.mu.Unlock()

	mask := uint32(1 << (position - 1))
	return oldNum&mask > 0, nil
}

func (that *ibitright) Set(position Position, oldNum uint32) (uint32, error) {
	if position < 1 || position > MaxPosition {
		return 0, ErrInvalidPosition
	}

	that.mu.Lock()
	defer that.mu.Unlock()

	return that.setBit(position, true, oldNum), nil
}

func (that *ibitright) UnSet(position Position, oldNum uint32) (uint32, error) {
	if position < 1 || position > MaxPosition {
		return 0, ErrInvalidPosition
	}

	that.mu.Lock()
	defer that.mu.Unlock()
	return that.setBit(position, false, oldNum), nil
}

func (that *ibitright) setBit(position Position, value bool, num uint32) uint32 {
	t := uint32(1 << (position - 1))
	if value {
		if num != 0 {
			t = num | t
		}
	} else {
		if num != 0 {
			t = num &^ t
		} else {
			t = 0
		}
	}
	return t
}
