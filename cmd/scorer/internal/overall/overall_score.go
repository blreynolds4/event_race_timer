package overall

import (
	"blreynolds4/event-race-timer/internal/competitors"
	"blreynolds4/event-race-timer/internal/results"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

const resultChunkSize = 25

type OverallScorer struct {
	logger         *slog.Logger
	overallResults []OverallResult
	rawResults     map[int]results.RaceResult // bib to result, keep latest result
}

type OverallResult struct {
	Athlete    *competitors.Competitor
	Finishtime time.Duration
	Place      int
	Bib        int
}

func NewOverallResults(l *slog.Logger) OverallScorer {
	return OverallScorer{
		logger:         l.With("scorer", "overall"),
		overallResults: make([]OverallResult, 0),
		rawResults:     make(map[int]results.RaceResult),
	}
}

func (ovr *OverallScorer) ScoreResults(ctx context.Context, source results.ResultStream) error {
	placeMap := make(map[int]OverallResult)

	// want to keep trying until told to stop via context
	results := make([]results.RaceResult, resultChunkSize)
	resultCount, err := source.GetResults(ctx, results)
	if err != nil {
		ovr.logger.Error("overall scorer error", "error", err)
		return fmt.Errorf("overall scorer error %w", err)
	}

	// get new results until the stream is empty
	for resultCount > 0 {
		// add any new results read to the raw storage
		for i := 0; i < resultCount; i++ {
			newResult := results[i]
			ovr.rawResults[newResult.Bib] = newResult
			ovr.logger.Debug("overall scorer adding result", "bib", newResult.Bib, "time", newResult.Time, "place", newResult.Place)
		}

		resultCount, err = source.GetResults(ctx, results)
		if err != nil {
			return err
		}
	}

	// no more results, so we can
	// build the output
	for _, result := range ovr.rawResults {
		placeMap[result.Place] = OverallResult{Athlete: result.Athlete, Place: result.Place, Finishtime: result.Time, Bib: result.Bib}
	}

	//output in the correct order
	ovr.overallResults = make([]OverallResult, 0)

	f, err := os.Create("overall_results.html")
	if err != nil {
		return err
	}
	defer f.Close()

	// write an html header just to the file
	_, err = f.WriteString("<html><head><title>Overall Results</title></head><body><pre>\n")
	if err != nil {
		return err
	}

	w := io.MultiWriter(f, os.Stdout)
	fmt.Printf("%s", "\x1Bc") // clear stdout
	fmt.Printf("Last Updated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "\n\n\n")
	fmt.Fprintln(w, "Place Bib   Name                             Grade Team                             Time")
	fmt.Fprintln(w, "===== ===== ================================ ===== ================================ ========")
	for i := 1; i <= len(placeMap); i++ {
		ovr.overallResults = append(ovr.overallResults, placeMap[i])
		r, exists := placeMap[i]
		if exists {
			fmt.Fprintf(w, "%-5d %-5d %-32s %-5d %-32s %-8s\n", r.Place, r.Bib, r.Athlete.Name, r.Athlete.Grade, r.Athlete.Team, formatFinishTime(r.Finishtime))
		}
	}

	_, err = f.WriteString("\n</pre></body></html>\n")
	if err != nil {
		return err
	}

	return nil
}

func formatFinishTime(t time.Duration) string {
	return time.Unix(0, 0).UTC().Add(time.Duration(t)).Format("04:05.00")
}
