package polling

import (
	"context"
	"errors"
	"testing"
	"time"
)

type processorFunc func(context.Context) (bool, error)

func (f processorFunc) ProcessNext(ctx context.Context) (bool, error) { return f(ctx) }

func TestWorkerContinuesAfterFailureAndStopsOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	calls := 0
	processor := processorFunc(func(context.Context) (bool, error) {
		calls++
		if calls == 1 {
			return true, errors.New("provider failure")
		}
		if calls == 2 {
			return false, nil
		}
		cancel()
		return true, nil
	})
	if err := Run(ctx, processor, time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Fatalf("unexpected processing count: %d", calls)
	}
}

func TestWorkerCancellationInterruptsIdleWait(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	called := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, processorFunc(func(context.Context) (bool, error) { close(called); return false, nil }), time.Hour)
	}()
	<-called
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("worker did not stop")
	}
}
