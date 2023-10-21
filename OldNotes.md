This assumes there are no RFID readers

Using a single computer (no local network required)

run redis in background

Create the athlete lookup in json or csv?
    Could also right code to take team names and default names to team_1, team_2, ...

run cli to generate start event
run cli to generate finish events  (this needs bib numbers to work)
run cli to generate place events (this needs bib numbers)
update cli to generate default place events from finish order?
    provisional place

run result builder to generate results:  overall and team scores saved to text file

Open Questions/Pending Changes:
What should goal be Tuesday if we don't have bibs?
Give out bibs and get names at the end?

How do we handle adding bibs to noBib finish events?
    need a view of finish events with no bibs.  The view provides a way to send finish events with correct bib
    Place events auto generated based on finish events sorted by duration.
        Ignore finish events with no bib
        Keep a cache of finish events (all with bibs) then re-issue place events when a new finish event changes the sort order

Manually generate place events that need to change based on chute ordering

==========================================================
Race Timing steps:
Start default placer with correct race name
Start scorer with correct config and race name/results name
Seed start/finish events with event generator (ensure correct source names in config)
Verify results


