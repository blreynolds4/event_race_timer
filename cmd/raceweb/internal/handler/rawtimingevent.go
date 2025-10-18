package handler

/*
JSON structure for OpenSignups timing event:

	{
	  "timestamp": 1677720000,  // (new) time in milliseconds since epoch when event happened
	    *** not sent "startTime": 1677720000, // time event started? or is it already when event happened?
	  "elapsedTime": 1000, // time in milliseconds since start time
	  "captureMode": "finish",
	  "antenna": 1,
	  "bib": "123", // chip/bib number
	  "host": "reader ip/hostname"
	}
*/
type OpenSignupsTimingEvent struct {
	EventTime   int    `json:"timestamp"`   // time the event was created
	ElapsedTime int    `json:"elapsedTime"` // time in milliseconds since start time set on RaceManager in AdjustStartTime dialog
	CaptureMode string `json:"captureMode"`
	// StartTime   int    `json:"startTime"`
	Antenna int    `json:"antenna"`
	Bib     string `json:"bib"`
	Host    string `json:"host"`
}
