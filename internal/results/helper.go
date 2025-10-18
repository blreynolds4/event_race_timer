package results

import (
	"blreynolds4/event-race-timer/internal/meets"
	"context"
	"time"
)

type MockResultStream struct {
	Get        func(ctx context.Context, result *meets.RaceResult, timeout time.Duration) (int, error)
	Results    []meets.RaceResult
	CancelFunc func()
}

func (mrs *MockResultStream) GetResult(ctx context.Context, result *meets.RaceResult, timeout time.Duration) (int, error) {
	if mrs.Get != nil {
		return mrs.Get(ctx, result, timeout)
	}

	if len(mrs.Results) > 0 {
		*result = mrs.Results[0]
		mrs.Results = mrs.Results[1:]
		return 1, nil
	}

	// nothing to return, cancel the context
	mrs.CancelFunc()
	return 0, nil
}
