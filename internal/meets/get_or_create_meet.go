package meets

import "io"

type MeetFinder interface {
	GetMeet(name string) (*Meet, error)
	io.Closer
}

type meetFinderImpl struct {
	reader MeetReader
	writer MeetWriter
}

func NewMeetFinder(rdr MeetReader, wrtr MeetWriter) (MeetFinder, error) {
	return &meetFinderImpl{
		reader: rdr,
		writer: wrtr,
	}, nil
}

func (m *meetFinderImpl) GetMeet(name string) (*Meet, error) {
	foundMeet, err := m.reader.GetMeet(name)
	if err != nil {
		return nil, err
	}
	if foundMeet != nil {
		// get the races for the meet
		races, err := m.reader.GetMeetRaces(foundMeet)
		if err != nil {
			return nil, err
		}
		foundMeet.races = races
		return foundMeet, nil
	}

	// Meet not found, create it
	newMeet := &Meet{Name: name}
	createdMeet, err := m.writer.SaveMeet(newMeet)
	if err != nil {
		return nil, err
	}
	return createdMeet, nil
}

func (m *meetFinderImpl) Close() error {
	err := m.reader.Close()
	if err != nil {
		return err
	}
	err = m.writer.Close()
	if err != nil {
		return err
	}
	return nil
}
