package command

import (
	"blreynolds4/event-race-timer/raceevents"
	"blreynolds4/event-race-timer/stream"
	"encoding/json"
)

func buildActualResults(rawOutput *stream.MockStream) []raceevents.Event {
	result := make([]raceevents.Event, len(rawOutput.Events))
	for i, msg := range rawOutput.Events {
		err := json.Unmarshal(msg.Data, &result[i])
		if err != nil {
			panic(err)
		}
	}

	return result
}

func buildEventMessages(testEvents []raceevents.Event) []stream.Message {
	result := make([]stream.Message, len(testEvents))
	for i, e := range testEvents {
		eData, err := json.Marshal(e)
		if err != nil {
			panic(err)
		}
		result[i] = stream.Message{
			ID:   e.ID,
			Data: eData,
		}
	}

	return result
}

func toMsg(r raceevents.Event) stream.Message {
	var msg stream.Message
	msgData, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	msg.Data = msgData

	return msg
}
