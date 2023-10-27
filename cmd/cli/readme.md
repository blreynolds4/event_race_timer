# Manual CLI Commands

## List Finishes missing Bibs and Add Bib
Need to get a list of finishes missing bib with:
list

Then select a message id to assign a bib to and use:
bib <finish msg id> <bib>

## Ping to be sure backend is alive
ping
Should respond with "pong"

## Set a place
place | p <bib> <place>

## Set a range of places by bib
placeRange | pr <next place>

At the "places" prompt provide bib to send place event and move to next place automatically
To exit the "places" prompt type "q" | "quit"

## Start a race
start <duration> | s <duration>
<duration> -> 10s or 5m or 2m30s (go duration syntax)

## Finish for a bib
f <bib> | finish <bib> 


## Exit the cli
q | quit | exit | stop

