package ierror

import (
	"log"
	"testing"
)

func TestError(t *testing.T) {
	if err := demo("ok"); err != nil {
		log.Println(err)
	}
	if err := demo("ok2"); err != nil {
		log.Println(err)
	}
}

func demo(str string) error {
	if str == "ok" {
		return nil
	} else {
		return New(44, "xzdemo")
	}

}
