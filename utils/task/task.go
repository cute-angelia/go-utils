package task

import (
	"github.com/robfig/cron/v3"
	"log"
)

type Task struct {
	Cron *cron.Cron
}

func NewTask() *Task {
	c := cron.New(cron.WithSeconds())
	c.Start()
	return &Task{
		Cron: c,
	}
}

func (self *Task) AddTask(spec string, f func()) {
	if id, err := self.Cron.AddFunc(spec, f); err != nil {
		log.Println("添加任务失败", err)
	} else {
		log.Println("task:", id)
	}
}

func (self *Task) StopTask() {
	self.Cron.Stop()
}
