package goshmuffle

import (
	"errors"
	"fmt"
	"syscall"
)

func (r *runner) Terminate() error {
	if r.exec == nil {
		return errors.New("exec is nil")
	}
	if r.cancel == nil {
		return errors.New("cancel func is nil")
	}
	defer r.cancel()

	// terminate the tree of processes
	pgid, err := syscall.Getpgid(r.exec.Process.Pid)
	if err != nil {
		return fmt.Errorf("get pid: %w", err)
	}

	return syscall.Kill(-pgid, 15)
}
