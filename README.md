# event_race_timer

A repo for event driving race timing

# manual_cli
Generate start and finish events to a redis stream based on event name.  These events can be refined along with events from other sources like chip readers to assign places and times to comptitors then provide scoring.

# event_api
Restful server for:
* recieving events (POST) from chip readers or other sources that can generate webhooks when they see a chip
* reading events from an event stream to support a scoring/timing UI
    read events live from stream
    re-read the whole stream again (use a new group name)

