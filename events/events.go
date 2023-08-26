package events

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

type EventType string

const (
	StartEventType  EventType = "StartEvent"
	FinishEventType EventType = "FinishEvent"
	PlaceEventType  EventType = "PlaceEvent"

	NoBib = -1
)

const (
	startTimeData  = "StartTime"
	finishTimeData = "FinishTime"
	bibData        = "Bib"
	placeData      = "Place"
)

type RaceEvent interface {
	GetID() string
	GetSource() string
	GetType() EventType
	GetTime() time.Time
}

// Start event will have type StartEvent and Data:
// StartTime
type StartEvent interface {
	RaceEvent
	GetStartTime() time.Time
}

// Finish event will have the type FinishEvent and Data:
// Bib and Finish Time
type FinishEvent interface {
	RaceEvent
	GetFinishTime() time.Time
	GetBib() int
}

type PlaceEvent interface {
	RaceEvent
	GetBib() int
	GetPlace() int
}

// all the data for all event types is the same underneath
// so all can be sent and read as Race Events
type EventTarget interface {
	SendRaceEvent(re RaceEvent) error
}

type EventSource interface {
	GetRaceEvent(timeout time.Duration) (RaceEvent, error)
}

// struct or interface?  what methods? enum for Data keys?
// Race event could be a public struct, then we can have interfaces for each
// event type that give access to info in Data via method?
type raceEvent struct {
	ID        string
	Source    string
	EventTime time.Time
	Type      EventType
	Data      map[string]any
}

func (re *raceEvent) GetID() string {
	return re.ID
}

func (re *raceEvent) GetTime() time.Time {
	return re.EventTime
}

func (re *raceEvent) GetType() EventType {
	return re.Type
}

func (re *raceEvent) GetSource() string {
	return re.Source
}

func NewStartEvent(src string, startTime time.Time) StartEvent {
	result := new(raceEvent)
	result.Data = make(map[string]interface{})
	result.Source = src
	result.EventTime = startTime
	result.Type = StartEventType

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

	// fmt.Println("Int TYPE:", fmt.Sprintf("%v", reflect.TypeOf(d)))
	result, ok := d.(int)
	// fmt.Println("result", result, ok)
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

func NewFinishEvent(src string, finishTime time.Time, bib int) FinishEvent {
	result := new(raceEvent)
	result.Data = make(map[string]interface{})
	result.Source = src
	result.EventTime = finishTime
	result.Type = FinishEventType

	// add the start time to the data payload
	result.Data[finishTimeData] = finishTime
	result.Data[bibData] = bib

	return result
}

func NewPlaceEvent(src string, bib, place int) PlaceEvent {
	result := new(raceEvent)
	result.Data = make(map[string]interface{})
	result.Source = src
	result.EventTime = time.Now().UTC()
	result.Type = PlaceEventType

	// add the start time to the data payload
	result.Data[placeData] = place
	result.Data[bibData] = bib

	return result
}

type redisEventStream struct {
	client    *redis.Client
	stream    string
	lastMsgId string
}

func NewRedisStreamEventTarget(c *redis.Client, name string) EventTarget {
	return &redisEventStream{
		client:    c,
		stream:    name,
		lastMsgId: "0",
	}
}

func NewRedisStreamEventSource(c *redis.Client, name string) EventSource {
	return &redisEventStream{
		client:    c,
		stream:    name,
		lastMsgId: "0",
	}
}

func (rset *redisEventStream) SendRaceEvent(re RaceEvent) error {
	// convert our event to a json to embed in the message
	eventData, err := json.Marshal(re)
	if err != nil {
		return err
	}

	addArgs := redis.XAddArgs{
		Stream: rset.stream,
		Values: map[string]interface{}{
			"event_type": string(re.GetType()),
			"event_time": re.GetTime().UnixMilli(),
			"source":     re.GetSource(),
			"event":      string(eventData),
		},
	}

	result := rset.client.XAdd(&addArgs)
	if result.Err() != nil {
		return result.Err()
	}

	fmt.Println("ok -", result.Val())
	return nil
}

func (rset *redisEventStream) GetRaceEvent(timeout time.Duration) (RaceEvent, error) {
	data, err := rset.client.XRead(&redis.XReadArgs{
		Streams: []string{rset.stream, rset.lastMsgId},
		//count is number of entries we want to read from redis
		Count: 1,
		//we use the block command to make sure if no entry is found we wait
		//timeout duration, 0 is forever
		Block: timeout,
	}).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err != redis.Nil && len(data[0].Messages) > 0 {
		msg := data[0].Messages[0]
		// create a result message and deserialize
		result := raceEvent{}
		data, ok := msg.Values["event"].(string)
		if ok {
			err = json.Unmarshal([]byte(data), &result)
			if err != nil {
				return nil, err
			}

			rset.lastMsgId = msg.ID
			result.ID = msg.ID
			return &result, nil
		}
	}
	return nil, nil
}

func (et EventType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", et)), nil
}

func (et *EventType) UnmarshalJSON(data []byte) error {
	x, _ := strconv.Unquote(string(data))
	*et = EventType(x)
	return nil
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
