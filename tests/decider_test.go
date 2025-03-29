package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/yandzee/paradox/retry"
)

var (
	ErrTest1   = errors.New("Test error 1")
	ErrTest2   = errors.New("Test error 2")
	ErrWrapped = errors.Join(ErrTest2, errors.New("Wrapped error"))
)

type DeciderTestDescriptor struct {
	Attempts         int
	Error            error
	Attempt          int
	DecideFn         func(error) retry.Decision
	ExpectedDecision retry.Decision
}

func TestDefaultDecider(t *testing.T) {
	runDeciderTests(t, []DeciderTestDescriptor{
		{
			Attempts:         0,
			Attempt:          0,
			Error:            nil,
			ExpectedDecision: retry.Success,
		},
		{
			Attempts:         0,
			Attempt:          100,
			Error:            nil,
			ExpectedDecision: retry.Success,
		},
		{
			Attempts:         1,
			Attempt:          0,
			Error:            nil,
			ExpectedDecision: retry.Success,
		},
		{
			Attempts:         1,
			Attempt:          1,
			Error:            nil,
			ExpectedDecision: retry.Success,
		},
		{
			Attempts:         1,
			Attempt:          2,
			Error:            nil,
			ExpectedDecision: retry.Success,
		},
		{
			Attempts:         0,
			Attempt:          0,
			Error:            ErrTest1,
			ExpectedDecision: retry.Retry,
		},
		{
			Attempts:         0,
			Attempt:          100,
			Error:            ErrTest1,
			ExpectedDecision: retry.Retry,
		},
		{
			Attempts:         5,
			Attempt:          1,
			Error:            ErrTest1,
			ExpectedDecision: retry.Retry,
		},
		{
			Attempts:         5,
			Attempt:          4,
			Error:            ErrTest1,
			ExpectedDecision: retry.Retry,
		},
		{
			Attempts:         5,
			Attempt:          5,
			Error:            ErrTest1,
			ExpectedDecision: retry.Fail,
		},
		{
			Attempts:         1,
			Attempt:          1,
			Error:            ErrTest1,
			ExpectedDecision: retry.Fail,
		},
		{
			Attempts:         1,
			Attempt:          2,
			Error:            ErrTest1,
			ExpectedDecision: retry.Fail,
		},
	})
}

func TestCustomDecider(t *testing.T) {
	decider := func(err error) retry.Decision {
		switch {
		case errors.Is(err, ErrTest1):
			return retry.Success
		case errors.Is(err, ErrWrapped):
			return retry.Retry
		case errors.Is(err, ErrTest2):
			return retry.Fail
		}

		return retry.Success
	}

	runDeciderTests(t, []DeciderTestDescriptor{
		{
			Attempts:         0,
			Attempt:          0,
			Error:            ErrTest1,
			DecideFn:         decider,
			ExpectedDecision: retry.Success,
		},
	})
}

func runDeciderTests(t *testing.T, descs []DeciderTestDescriptor) {
	for i, td := range descs {
		t.Run(
			fmt.Sprintf("Decider test %d", i),
			func(t *testing.T) {
				decider := &retry.CustomDecider{
					Attempts:    td.Attempts,
					DecideError: td.DecideFn,
				}

				rctx := &retry.RetryContext{
					Attempt:   td.Attempt,
					LastError: td.Error,
				}

				got := decider.Decide(rctx)
				if got != td.ExpectedDecision {
					t.Fatalf("Decision: got %v, expected: %v", got, td.ExpectedDecision)
				}
			},
		)
	}
}
