package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryFailed(t *testing.T) {
	ast := assert.New(t)
	var number int
	increaseNumber := func() error {
		number++
		return errors.New("error occurs")
	}

	err := Retry(increaseNumber, RetryDuration(time.Microsecond*50))

	ast.Error(err)
	ast.Equal(DefaultRetryTimes, number, "TestRetryFailed")
}

func TestRetrySucceeded(t *testing.T) {
	ast := assert.New(t)

	var number int
	increaseNumber := func() error {
		number++
		if number == DefaultRetryTimes {
			return nil
		}
		return errors.New("error occurs")
	}

	err := Retry(increaseNumber, RetryDuration(time.Microsecond*50))

	ast.Nil(err, "TestRetrySucceeded")
	ast.Equal(DefaultRetryTimes, number, "TestRetrySucceeded")
}

func TestSetRetryTimes(t *testing.T) {
	ast := assert.New(t)

	var number int
	increaseNumber := func() error {
		number++
		return errors.New("error occurs")
	}

	err := Retry(increaseNumber, RetryDuration(time.Microsecond*50), RetryTimes(3))

	ast.Error(err)
	ast.Equal(3, number, TestSetRetryTimes)
}

func TestCancelRetry(t *testing.T) {
	ast := assert.New(t)

	ctx, cancel := context.WithCancel(context.TODO())
	var number int
	increaseNumber := func() error {
		number++
		if number > 3 {
			cancel()
		}
		return errors.New("error occurs")
	}

	err := Retry(increaseNumber,
		RetryDuration(time.Microsecond*50),
		Context(ctx),
	)

	ast.Error(err)
	ast.Equal(4, number, "TestCancelRetry")
}
