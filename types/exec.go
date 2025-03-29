package types

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"
)

type exec struct {
	name          string
	status        Status
	fn            func(io.Writer, *exec) error
	fnWithContext func(context.Context, io.Writer, *exec) error
	out           bytes.Buffer
	err           error
	ctx           context.Context
	cancel        context.CancelFunc

	TaskDone  int
	TaskTotal int
}

func NewExec(name string, fn func(io.Writer, *exec) error) *exec {
	ctx, cancel := context.WithCancel(context.Background())
	return &exec{name, PENDING, fn, nil, bytes.Buffer{}, nil, ctx, cancel, 0, 0}
}

func NewExecWithContext(name string, fn func(context.Context, io.Writer, *exec) error) *exec {
	ctx, cancel := context.WithCancel(context.Background())
	return &exec{name, PENDING, nil, fn, bytes.Buffer{}, nil, ctx, cancel, 0, 0}
}

func (e *exec) Description() string {
	return e.name
}

func (e *exec) Run() {
	go func() {
		finished := make(chan any)

		go func() {
			e.status = RUNNING

			if e.fn != nil {
				e.err = e.fn(&e.out, e)
			} else if e.fnWithContext != nil {
				e.err = e.fnWithContext(e.ctx, &e.out, e)
			} else {
				e.err = errors.New("no function defined")
			}

			if e.ctx.Err() != nil {
				e.status = CANCELLED
			} else if e.err != nil {
				e.status = ERROR
			} else {
				e.status = DONE
			}
			close(finished)
		}()

		select {
		case <-e.ctx.Done():
			time.Sleep(50 * time.Millisecond) // time for fnWithContext to finish itself
			if e.status == PENDING || e.status == RUNNING {
				e.status = CANCELLED
			}

		case <-finished:
			// Nothing
		}
	}() // inner goroutine will be killed at this point
}

func (e *exec) Cancel() {
	e.cancel()
}

func (e *exec) Status() Status {
	return e.status
}

func (e *exec) Progress() (int, int) {
	return e.TaskDone, e.TaskTotal
}

func (e *exec) Log() string {
	return e.out.String()
}

func (e *exec) Err() error {
	return e.err
}
