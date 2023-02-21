package random

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRandString(t *testing.T) {
	ast := assert.New(t)

	pattern := `^[a-zA-Z]+$`
	reg := regexp.MustCompile(pattern)

	randStr := RandString(6, LetterAbc)

	ast.Equal(6, len(randStr))
	ast.Equal(true, reg.MatchString(randStr))
}
