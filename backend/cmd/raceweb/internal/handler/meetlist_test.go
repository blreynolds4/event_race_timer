package handler

import (
	"blreynolds4/event-race-timer/internal/meets"
	"log/slog"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// GetMeet(name string) (*Meet, error)
// GetMeets() ([]*Meet, error)
// GetMeetRaces(m *Meet) ([]Race, error)
// io.Closer

type MockMeetReader struct {
	GetMeetByNameFunc func(name string) (*meets.Meet, error)
	GetAllMeetsFunc   func() ([]*meets.Meet, error)
	GetMeetRacesFunc  func(m *meets.Meet) ([]meets.Race, error)
	CloseFunc         func() error
}

func (m *MockMeetReader) GetMeet(name string) (*meets.Meet, error) {
	if m.GetMeetByNameFunc != nil {
		return m.GetMeetByNameFunc(name)
	}
	return nil, nil
}

func (m *MockMeetReader) GetMeets() ([]*meets.Meet, error) {
	if m.GetAllMeetsFunc != nil {
		return m.GetAllMeetsFunc()
	}
	return nil, nil
}

func (m *MockMeetReader) GetMeetRaces(meet *meets.Meet) ([]meets.Race, error) {
	if m.GetMeetRacesFunc != nil {
		return m.GetMeetRacesFunc(meet)
	}
	return nil, nil
}

func (m *MockMeetReader) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func TestMeetListHandler(t *testing.T) {
	// Create a logger with the TestHandler
	logger := slog.New(slog.DiscardHandler)

	req := httptest.NewRequest("GET", "/api/meets", nil)
	// Create a new Gin router
	router := gin.Default()

	mockMeetReader := &MockMeetReader{
		GetAllMeetsFunc: func() ([]*meets.Meet, error) {
			return []*meets.Meet{
				{Name: "Meet 1"},
				{Name: "Meet 2"},
			}, nil
		},
	}

	// Register the handler
	router.GET("/api/meets", NewMeetListHandler(mockMeetReader, logger))

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := `[{"Name":"Meet 1"},{"Name":"Meet 2"}]`
	assert.JSONEq(t, expectedBody, w.Body.String())
}
