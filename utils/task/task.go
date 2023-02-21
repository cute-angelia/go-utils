package task

import (
	"github.com/robfig/cron/v3"
	"log"
)

/*
Entry                  | Description                                | Equivalent To
-----                  | -----------                                | -------------
@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
@monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
@weekly                | Run once a week, midnight between Sat/Sun  | 0 0 0 * * 0
@daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
@hourly                | Run once an hour, beginning of hour        | 0 0 * * * *
*/

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
