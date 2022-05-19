package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToUnix(t *testing.T) {
	assert := assert.New(t)
	tm1 := NewUnixNow()
	unixTimestamp := tm1.GetUnix()

	tm2 := NewUnix(unixTimestamp)

	assert.Equal(tm1, tm2, "TestToUnix")

	//  获取0点
	t.Log(tm2.GetTimeZero().GetUnix())
	t.Log(tm2.GetTimeZero().FormatDateTime())
}
