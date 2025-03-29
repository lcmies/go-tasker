package types

import "time"

type Task interface {
	Description() string
	Run()
	Cancel()
	Status() Status
	Progress() (int, int)
	Log() string
	Err() error
}

type taskForMgr struct {
	Task
	added time.Time
}

type TaskDetail struct {
	Id        uint64    `json:"id"`
	Added     time.Time `json:"added"`
	Name      string    `json:"name"`
	StatusStr string    `json:"status"`
	Log       string    `json:"log,omitempty"`
	ErrStr    string    `json:"err,omitempty"`
	TaskDone  int       `json:"task_done"`
	TaskTotal int       `json:"task_total"`
}
