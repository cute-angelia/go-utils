package task

import (
	"log"
	"testing"
)

func TestFirstDo(t *testing.T) {
	itask := NewTask()
	itask.AddTask("3 * * * * *", func() {
		log.Println("i here")
	})
	defer itask.StopTask()
	select {

	}
}
