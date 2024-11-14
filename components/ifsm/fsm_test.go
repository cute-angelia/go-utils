package ifsm

import (
	"context"
	"log"
	"testing"

	"github.com/looplab/fsm"
)

const (
	FsmStatusOpen   string = "open"
	FsmStatusClosed string = "closed"
	FsmStatusRun    string = "run"
)

type Door struct {
	Name string
	FSM  *fsm.FSM
}

func NewDoor(name string) *Door {
	d := &Door{
		Name: name,
	}

	d.FSM = fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "goOpen", Src: []string{FsmStatusClosed}, Dst: FsmStatusOpen},
			{Name: "goRun", Src: []string{FsmStatusClosed}, Dst: FsmStatusRun},
			{Name: "goClose", Src: []string{FsmStatusOpen}, Dst: FsmStatusClosed},
		},
		fsm.Callbacks{
			"enter_state": func(_ context.Context, e *fsm.Event) { d.newState(e) },
		},
	)

	return d
}

func (d *Door) newState(e *fsm.Event) {
	log.Printf(">> 1 newState, current state:%s, Dst:%s \n", d.FSM.Current(), e.Dst)
}

func TestFms(t *testing.T) {

	log.SetFlags(log.Lshortfile)

	door := NewDoor("测试")

	log.Printf("init fsm current state: %s \n", door.FSM.Current())

	err := door.FSM.Event(context.Background(), "goOpen")
	if err != nil {
		log.Println(err)
	}
	log.Printf("fsm current state: %s \n", door.FSM.Current())

	// open -> run 不允许
	err = door.FSM.Event(context.Background(), "goRun")
	if err != nil {
		log.Println("error:  open -> run 不允许", err)
	}
	log.Printf("fsm current state: %s \n", door.FSM.Current())

	err = door.FSM.Event(context.Background(), "goClose")
	if err != nil {
		log.Println(err)
	}
	log.Printf("fsm current state: %s \n", door.FSM.Current())

	err = door.FSM.Event(context.Background(), "goRun")
	if err != nil {
		log.Println(err)
	}
	log.Printf("fsm current state: %s \n", door.FSM.Current())

	err = door.FSM.Event(context.Background(), "goClose")
	if err != nil {
		log.Println(err)
	}
	log.Printf("fsm current state: %s \n", door.FSM.Current())

	err = door.FSM.Event(context.Background(), "goClose")
	if err != nil {
		log.Println(err)
	}
	log.Printf("fsm current state: %s \n", door.FSM.Current())
}
