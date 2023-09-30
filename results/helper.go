package results

import (
	"context"
	"time"
)

type MockResultSource struct {
	Get        func(ctx context.Context, result *RaceResult, timeout time.Duration) (int, error)
	Results    []RaceResult
	CancelFunc func()
}

func (mrs *MockResultSource) GetResult(ctx context.Context, result *RaceResult, timeout time.Duration) (int, error) {
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
