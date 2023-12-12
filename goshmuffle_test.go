package goshmuffle_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/pershinov/goshmuffle"
)

type res struct {
	s string
}

func (r *res) Store(s string) {
	r.s = s
}

func TestRun(t *testing.T) {
	expected := "hello"
	r := &res{}
	cmd := goshmuffle.New("echo", expected).WithOut(r)

	err := cmd.Run(context.Background())
	if err != nil {
		t.Errorf("run: %s", err.Error())
	}

	if r.s != expected {
		t.Errorf("out: %s is not expected: %s", r.s, expected)
	}
}

func TestTerminate(t *testing.T) {
	cmd := goshmuffle.New("sleep", "10")
	go func() {
		err := cmd.Run(context.Background())
		if err != nil && !strings.Contains(err.Error(), "terminated") {
			t.Errorf("run: %s", err.Error())
		}
	}()

	for !cmd.IsRunning() {
	}
	timestamp := time.Now()

	if cmd.IsDone() {
		t.Errorf("cmd cannot be done")
	}

	err := cmd.Terminate()
	if err != nil {
		t.Errorf("terminate: %s", err.Error())
	}

	for cmd.IsRunning() {
	}

	if !cmd.IsDone() {
		t.Errorf("cmd cannot be undone")
	}

	execTime := time.Now().Sub(timestamp)
	limit := 2 * time.Second
	if execTime > limit {
		t.Errorf("time limit %s is exceeded with %s", limit.String(), execTime.String())
	}
}
