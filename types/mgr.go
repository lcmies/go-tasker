package types

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Mgr struct {
	tasks map[uint64]*taskForMgr
	lock  sync.Mutex
}

func NewMgr() *Mgr {
	return &Mgr{
		tasks: make(map[uint64]*taskForMgr),
	}
}

func (m *Mgr) Clear() {
	for id, task := range m.tasks {
		if task.Status() == ERROR || task.Status() == CANCELLED || task.Status() == DONE {
			delete(m.tasks, id)
		}
	}
}

func (m *Mgr) Cancel(id uint64) {
	task, ok := m.tasks[id]
	if !ok {
		return
	}
	if task.Status() == PENDING || task.Status() == RUNNING {
		task.Cancel()
	}
}

func (m *Mgr) Add(task Task) {
	m.lock.Lock()
	defer m.lock.Unlock()

	var id uint64
	for {
		id = uint64(rand.Int31())
		if _, ok := m.tasks[id]; !ok {
			break
		}
	}
	m.tasks[id] = &taskForMgr{task, time.Now()}
	task.Run()
}

func (m *Mgr) List() []TaskDetail {
	ret := make([]TaskDetail, 0, len(m.tasks))

	for id, task := range m.tasks {
		td := TaskDetail{
			Id:    id,
			Name:  task.Description(),
			Added: task.added,
			Log:   task.Log(),
		}

		td.TaskDone, td.TaskTotal = task.Progress()
		td.StatusStr = task.Status().Str()

		if task.Err() != nil {
			td.ErrStr = fmt.Sprintf("%+v", task.Err())
		}

		ret = append(ret, td)
	}

	return ret
}
