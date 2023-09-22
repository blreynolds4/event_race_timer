package eventstream

import (
	"blreynolds4/event-race-timer/events"
	"blreynolds4/event-race-timer/stream"
	"encoding/json"
	"fmt"
	"time"
)

const (
	startTimeData  = "StartTime"
	finishTimeData = "FinishTime"
	bibData        = "Bib"
	placeData      = "Place"
)

// struct or interface?  what methods? enum for Data keys?
// Race event could be a public struct, then we can have interfaces for each
// event type that give access to info in Data via method?
type raceEvent struct {
	ID        string
	Source    string
	EventTime time.Time
	Type      events.EventType
	Data      map[string]any
}

func RaceEventToStreamMessage(re events.RaceEvent) (stream.Message, error) {
	lre, ok := re.(*raceEvent)
	if !ok {
		return stream.Message{}, fmt.Errorf("not a supported event structure")
	}

	// convert our event to a json to embed in the message
	eventData, err := json.Marshal(lre)
	if err != nil {
		return stream.Message{}, err
	}

	msg := stream.Message{
		Values: map[string]interface{}{
			"event_type": string(re.GetType()),
			"event_time": re.GetTime().UnixMilli(),
			"source":     re.GetSource(),
			"event":      string(eventData),
		},
	}

	return msg, nil
}

func StreamMessageToRaceEvent(msg stream.Message) (events.RaceEvent, error) {
	re := raceEvent{
		ID: msg.ID,
	}

	data, ok := msg.Values["event"].(string)
	if ok {
		err := json.Unmarshal([]byte(data), &re)
		if err != nil {
			return nil, err
		}

		return &re, nil
	}

	return nil, fmt.Errorf("Values data was not a string, can't build RaceEvent")
}

func (re *raceEvent) GetID() string {
	return re.ID
}

func (re *raceEvent) GetTime() time.Time {
	return re.EventTime
}

func (re *raceEvent) GetType() events.EventType {
	return re.Type
}

func (re *raceEvent) GetSource() string {
	return re.Source
}

func NewStartEvent(src string, startTime time.Time) events.StartEvent {
	result := new(raceEvent)
	result.Data = make(map[string]interface{})
	result.Source = src
	result.EventTime = startTime
	result.Type = events.StartEventType

	// add the start time to the data payload
	result.Data[startTimeData] = startTime

	return result
}

func (re *raceEvent) GetStartTime() time.Time {
	return re.getTimeData(startTimeData)
}

func (re *raceEvent) getTimeData(field string) time.Time {
	timeData, found := re.Data[field]
	if !found {
		panic(fmt.Sprintf("data for event field %s is missing", field))
	}

	result, ok := timeData.(time.Time)
	if !ok {
		panic(fmt.Sprintf("%s data in event should be time.Time", field))
	}

	return result
}

func (re *raceEvent) GetFinishTime() time.Time {
	return re.getTimeData(finishTimeData)
}

func (re *raceEvent) getIntData(field string) int {
	d, found := re.Data[field]
	if !found {
		panic(fmt.Sprintf("data for event field %s is missing", field))
	}

	result, ok := d.(int)
	if !ok {
		panic(fmt.Sprintf("%s data in event should be an int", field))
	}

	return result
}

func (re *raceEvent) GetBib() int {
	return re.getIntData(bibData)
}

func (re *raceEvent) GetPlace() int {
	return re.getIntData(placeData)
}

func NewFinishEvent(src string, finishTime time.Time, bib int) events.FinishEvent {
	result := raceEvent{
		Data:      make(map[string]interface{}),
		Source:    src,
		EventTime: finishTime,
		Type:      events.FinishEventType,
	}

	// add the start time to the data payload
	result.Data[finishTimeData] = finishTime
	result.Data[bibData] = bib

	return &result
}

func NewPlaceEvent(src string, bib, place int) events.PlaceEvent {
	result := new(raceEvent)
	result.Data = make(map[string]interface{})
	result.Source = src
	result.EventTime = time.Now().UTC()
	result.Type = events.PlaceEventType

	// add the start time to the data payload
	result.Data[placeData] = place
	result.Data[bibData] = bib

	return result
}

func (et *raceEvent) UnmarshalJSON(data []byte) error {
	var objmap map[string]*json.RawMessage

	err := json.Unmarshal(data, &objmap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objmap["Source"], &et.Source)
	if err != nil {
		return err
	}
	err = json.Unmarshal(*objmap["EventTime"], &et.EventTime)
	if err != nil {
		return err
	}
	err = json.Unmarshal(*objmap["Type"], &et.Type)
	if err != nil {
		return err
	}
	err = json.Unmarshal(*objmap["Data"], &et.Data)
	if err != nil {
		return err
	}
	// want to convert all numbers top int
	// for each k,v convert to int if it's type is float64
	for k, v := range et.Data {
		if tempFloat, ok := v.(float64); ok {
			et.Data[k] = int(tempFloat)
		} else {
			et.Data[k] = v
		}
	}

	// convert the dates stored in the Data
	if et.Data[startTimeData] != nil {
		et.Data[startTimeData] = et.unmarshallDateData(et.Data[startTimeData])
	}

	if et.Data[finishTimeData] != nil {
		et.Data[finishTimeData] = et.unmarshallDateData(et.Data[finishTimeData])
	}

	return nil
}

func (et *raceEvent) unmarshallDateData(rawDate any) time.Time {
	result := time.Time{}
	data, ok := rawDate.(string)
	if !ok {
		panic(fmt.Sprintf("expected raw date to be string but was %v", rawDate))
	}

	result, err := time.Parse(time.RFC3339Nano, data)
	if err != nil {
		panic(err)
	}
	return result
}
