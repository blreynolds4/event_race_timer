package meets

type MockResultWriter struct {
	SavedResults []RaceResult
}

func NewMockResultWriter() *MockResultWriter {
	return &MockResultWriter{
		SavedResults: make([]RaceResult, 0),
	}
}

func (mrw *MockResultWriter) SaveResult(rr *RaceResult) (*RaceResult, error) {
	mrw.SavedResults = append(mrw.SavedResults, *rr)
	return rr, nil
}

func (mrw *MockResultWriter) Close() error {
	return nil
}
