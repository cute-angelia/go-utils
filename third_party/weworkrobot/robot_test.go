package weworkrobot

import "testing"

func TestSendText(t *testing.T) {
	cli := Load("b5140d8b-xxxx-4e75-bd29-e52e203a1090").Build(WithDebug(true))
	cli.SendText("good")
}
