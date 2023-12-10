package goshmuffle

import (
	"bufio"
	"errors"
)

type out interface {
	Store(s string)
}

func (r *runner) store() error {
	if r.stdout == nil {
		return errors.New("stdout is nil")
	}

	scanner := bufio.NewScanner(r.stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		r.out.Store(m)
	}

	return nil
}
