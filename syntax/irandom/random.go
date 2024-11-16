package random

import (
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
