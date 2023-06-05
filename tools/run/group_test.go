package run

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestZero(t *testing.T) {
	t.Parallel()

	var g Group
	res := make(chan error)
	go func() { res <- g.Run() }()
	select {
	case err := <-res:
		if err != nil {
			t.Errorf("%v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout")
	}
}

func TestOne(t *testing.T) {
	t.Parallel()

	myError := errors.New("foobar")
	var g Group
	g.AddNamed("", func() error { return myError }, func(error) {})
	res := make(chan error)
	go func() { res <- g.Run() }()
	select {
	case err := <-res:

		if want, have := myError, err; !strings.Contains(have.Error(), want.Error()) {
			t.Errorf("want %v, have %v", want, have)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout")
	}
}

func TestMany(t *testing.T) {
	t.Parallel()

	interrupt := errors.New("interrupt")
	var g Group
	g.AddNamed("", func() error { return interrupt }, func(error) {})
	cancel := make(chan struct{})
	g.AddNamed("", func() error { <-cancel; return nil }, func(error) { close(cancel) })
	res := make(chan error)
	go func() { res <- g.Run() }()
	select {
	case err := <-res:
		if want, have := interrupt, err; !strings.Contains(have.Error(), want.Error()) {
			t.Errorf("want %v, have %v", want, have)
		}
	case <-time.After(100 * time.Millisecond):
		t.Errorf("timeout")
	}
}

func TestHalted(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	g := NewNamedGroup(ctx, "test")

	long := func(ctx context.Context) error { time.Sleep(1100 * time.Millisecond); return nil }
	short := func(ctx context.Context) error { return nil }

	g.AddWithContextNamed("Long", long)
	g.AddWithContextNamed("Short", short)

	go cancel()

	err := g.Run()
	require.Error(t, err)
}
