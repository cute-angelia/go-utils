package logger

import (
	"log"
	"os"
)

func Example() {
	Info("hello, world.")
}

func ExampleNewLogger() {
	w := os.Stdout
	flag := log.Llongfile
	l := NewWriterLogger(w, flag, 3)
	l.Info("hello, world")
}
