package results

import (
	"context"
	"fmt"
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
		fmt.Println(len(mrs.Results))
		mrs.Results = mrs.Results[1:]
		fmt.Println(len(mrs.Results))
		return result, nil
	}

	return RaceResult{}, nil
}
