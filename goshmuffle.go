package goshmuffle

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

type result interface {
	Store(s string)
}

type runner struct {
	cmd  string
	args []string
	res  result

	exec   *exec.Cmd
	stderr io.ReadCloser
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

func (cr *runner) WithResult(r result) *runner {
	cr.res = r
	return cr
}

func (cr *runner) Run(ctx context.Context) (err error) {
	cr.exec = exec.CommandContext(ctx, cr.cmd, cr.args...)
	cr.exec.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cr.stderr, err = cr.exec.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	err = cr.exec.Start()
	if err != nil {
		return fmt.Errorf("start cmd: %w", err)
	}

	if cr.res != nil {
		err = cr.store()
		if err != nil {
			return fmt.Errorf("store result: %w", err)
		}
	}

	return cr.exec.Wait()
}

func (cr *runner) store() error {
	scanner := bufio.NewScanner(cr.stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		cr.res.Store(m)
	}

	return nil
}

func (cr *runner) Terminate() error {
	if cr.exec == nil {
		return errors.New("kill: exec is nil")
	}

	pgid, err := syscall.Getpgid(cr.exec.Process.Pid)
	if err != nil {
		return fmt.Errorf("get pid: %w", err)
	}

	return syscall.Kill(-pgid, 15)
}
