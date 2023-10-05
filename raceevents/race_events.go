package raceevents

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	ID        string
	EventTime time.Time
	Data      any
}

type StartEvent struct {
	Source    string
	StartTime time.Time
}

type FinishEvent struct {
	Source     string
	Bib        int
	FinishTime time.Time
}

type PlaceEvent struct {
	Source string
	Bib    int
	Place  int
}

type marshalledEvent struct {
	ID        string
	EventTime time.Time
	Data      any
	DataType  string `json:"dataType"`
}

func (e Event) MarshalJSON() ([]byte, error) {
	actual := marshalledEvent{
		ID:        e.ID,
		EventTime: e.EventTime,
		Data:      e.Data,
	}

	switch e.Data.(type) {
	case StartEvent:
		actual.DataType = "start"
	case FinishEvent:
		actual.DataType = "finish"
	case PlaceEvent:
		actual.DataType = "place"
	default:
		return nil, fmt.Errorf("unknown type in Event Data")
	}

	return json.Marshal(actual)
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(data, &objmap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objmap["ID"], &e.ID)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objmap["EventTime"], &e.EventTime)
	if err != nil {
		return err
	}

	var dataType string
	err = json.Unmarshal(*objmap["dataType"], &dataType)
	if err != nil {
		return err
	}

	switch dataType {
	case "start":
		var se StartEvent
		err = json.Unmarshal(*objmap["Data"], &se)
		e.Data = se
	case "finish":
		var fe FinishEvent
		err = json.Unmarshal(*objmap["Data"], &fe)
		e.Data = fe
	case "place":
		var pe PlaceEvent
		err = json.Unmarshal(*objmap["Data"], &pe)
		e.Data = pe
	default:
		return fmt.Errorf("unknown type in Event Data")
	}

	return nil
}
