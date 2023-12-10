package goshmuffle

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"
)

type runner struct {
	// config
	cmd  string
	args []string
	out  out

	// cmd
	exec   *exec.Cmd
	stdout io.ReadCloser

	// state
	state state
	mu    sync.RWMutex

	// context
	cancel context.CancelFunc
}

func New(
	cmd string,
	args ...string,
) *runner {
	return &runner{
		cmd:  cmd,
		args: args,
	}
}

func (r *runner) WithOut(o out) *runner {
	r.out = o
	return r
}

func (r *runner) Run(ctx context.Context) (err error) {
	ctx, r.cancel = context.WithCancel(ctx)

	// configure cmd
	r.exec = exec.CommandContext(ctx, r.cmd, r.args...)
	r.exec.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// configure stdout
	if r.out != nil {
		r.stdout, err = r.exec.StdoutPipe()
		if err != nil {
			return fmt.Errorf("stdout pipe: %w", err)
		}
	}

	// run cmd
	err = r.exec.Start()
	if err != nil {
		return fmt.Errorf("start cmd: %w", err)
	}

	// states control
	r.setStateRunning()
	defer func() {
		r.setStateDone()
	}()

	// store out
	if r.out != nil {
		err = r.store()
		if err != nil {
			return fmt.Errorf("store out: %w", err)
		}
	}

	return r.exec.Wait()
}

func (r *runner) IsRunning() bool {
	return r.isRunning()
}

func (r *runner) IsDone() bool {
	return r.isDone()
}
