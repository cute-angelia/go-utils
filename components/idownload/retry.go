package idownload

import (
	"context"
	"errors"
	"log"
	"time"
)

var (
	ErrRetryFail = errors.New("retry fail")
	ErrRetry     = errors.New("need to retry")
)

type Retry struct {
	attempt     int // Maximum number of attempts
	currAttempt int
	waitTime    time.Duration
	cb          func() error
}

func NewRetry(attempt int, waitTime time.Duration) *Retry {
	return &Retry{
		attempt:  attempt,
		waitTime: waitTime,
	}
}

func (r *Retry) Func(cb func() error) *Retry {
	r.cb = cb
	return r
}

func (r *Retry) reset() {
	r.currAttempt = 0
}

// Do send function
func (r *Retry) Do(ctx context.Context) (err error) {
	defer r.reset()

	tk := time.NewTimer(r.waitTime)
	defer tk.Stop()

	for i := 0; i < r.attempt; i++ {
		// 这里只要调用Func方法，且回调函数返回ErrRetry 会生成新的*http.Request对象
		// 不使用DataFlow.Do()方法原因基于两方面考虑
		// 1.为了效率只需经过一次编码器得到*http.Request,如果需要重试几次后面是多次使用解码器.Bind()函数
		// 2.为了更灵活的控制
		if r.cb != nil {
			if err = r.cb(); err != nil {
				if err == ErrRetry {
					log.Println("Retry:次数", i+1)
				} else {
					log.Println("Retry:失败，不继续下载", err)
				}
			} else {
				return nil
			}
		}

		tk.Reset(r.waitTime)
		select {
		case <-tk.C:
			// 外部可以使用context直接取消
		case <-ctx.Done():
			log.Println("Retry:重试超时取消")
			return ctx.Err()
		}
		r.currAttempt++
	}
	return ErrRetryFail
}
