package racearchive

import "blreynolds4/event-race-timer/internal/raceevents"

// don't incluce race name is the output, that allows
// un-archive to use the filename as the racename stream
type RaceArchive struct {
	RaceEvents []raceevents.Event
}
