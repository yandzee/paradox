package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/yandzee/paradox/backoff"
	"github.com/yandzee/paradox/clock"
	"github.com/yandzee/paradox/retry"
)

var ErrTest = errors.New("test err")

var TestBackoff = backoff.ConstantBackoff{
	Delay: 1 * time.Second,
}

var TestDecider = retry.NewAttemptsDecider(retry.InfiniteAttempts)

type TestDescriptor struct {
	DoReturns []DoReturn
	Error     error
	Time      int64
	Backoff   backoff.Backoff
	Decider   retry.Decider
}

type DoReturn struct {
	time.Duration
	error
}

func TestManualDurationSet(t *testing.T) {
	runTests(t, []TestDescriptor{
		{
			DoReturns: []DoReturn{
				{0, nil},
			},
			Time: int64(0 * time.Second),
		},
		{
			DoReturns: []DoReturn{
				{2 * time.Second, ErrTest},
				{0, nil},
			},
			Time: int64(2 * time.Second),
		},
		{
			DoReturns: []DoReturn{
				{0, ErrTest},
				{2 * time.Second, ErrTest},
				{0, nil},
			},
			Time: int64(3 * time.Second),
		},
		{
			DoReturns: []DoReturn{
				{0, ErrTest},
				{0, ErrTest},
				{0, ErrTest},
				{0, ErrTest},
				{0, ErrTest},
				{2 * time.Second, ErrTest},
				{0, nil},
			},
			Time: int64(7 * time.Second),
		},
		{
			DoReturns: []DoReturn{
				{0, ErrTest},
				{0, ErrTest},
				{0, ErrTest},
				{0, ErrTest},
				{0, ErrTest},
				{0, nil},
			},
			Time: int64(5 * time.Second),
		},
	})
}

func runTests(t *testing.T, tds []TestDescriptor) {
	for i, td := range tds {
		if td.Backoff == nil {
			td.Backoff = &TestBackoff
		}

		if td.Decider == nil {
			td.Decider = TestDecider
		}

		cl := &clock.MockClock{
			AutoFastForward: true,
		}

		r := retry.NewWithClock(td.Backoff, td.Decider, cl)
		fn := td.BuildDoFunction(t)

		err := r.DoManual(t.Context(), fn)

		if !errors.Is(err, td.Error) {
			t.Fatalf("Test %d: errors mismatch, got %v, expected %v", i, err, td.Error)
		}

		if cl.Current != td.Time {
			t.Fatalf(
				"Test %d: clock time mismatch, got %v (%s), expected %v (%s)",
				i,
				cl.Current,
				time.Duration(cl.Current).String(),
				td.Time,
				time.Duration(td.Time).String(),
			)
		}
	}
}

func (td *TestDescriptor) BuildDoFunction(t *testing.T) retry.DoManualFn {
	return func(rc *retry.RetryContext) (time.Duration, error) {
		if rc.Attempt >= len(td.DoReturns) {
			t.Fatalf("BuildDoFunction: attempt %d has no return pair", rc.Attempt)
		}

		pair := td.DoReturns[rc.Attempt]
		return pair.Duration, pair.error
	}
}
