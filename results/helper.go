package results

import (
	"context"
)

type MockResultSource struct {
	Get     func(ctx context.Context) (RaceResult, error)
	Results []RaceResult
}

func (mrs *MockResultSource) GetResult(ctx context.Context) (RaceResult, error) {
	if mrs.Get != nil {
		return mrs.Get(ctx)
	}

	if len(mrs.Results) > 0 {
		result := mrs.Results[0]
		mrs.Results = mrs.Results[1:]
		return result, nil
	}

	return RaceResult{}, nil
}
