package overall

import (
	"blreynolds4/event-race-timer/internal/meets"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

type OverallRaceScorer struct {
	race   *meets.Race
	logger *slog.Logger
}

func NewOverallRaceResults(race *meets.Race, l *slog.Logger) OverallRaceScorer {
	return OverallRaceScorer{
		race:   race,
		logger: l.With("scorer", "overall"),
	}
}

func (ovr *OverallRaceScorer) ScoreResults(ctx context.Context, resultsReader meets.RaceResultReader) error {
	// get results for the race the resultsReader was created for.
	// results are returned in place order
	ovr.logger.Info("Building overall...")
	raceResults, err := resultsReader.GetRaceResults()
	if err != nil {
		ovr.logger.Error("overall race scorer error", "error", err)
		return fmt.Errorf("overall race scorer error %w", err)
	}

	overallResults := make([]OverallResult, len(raceResults))

	// build the output
	for i, result := range raceResults {
		overallResults[i] = OverallResult{Athlete: result.Athlete, Place: result.Place, Finishtime: result.Time, Bib: result.Bib}
	}

	f, err := os.Create("overall_results.html")
	if err != nil {
		return err
	}
	defer f.Close()

	// write an html header just to the file
	_, err = f.WriteString("<html><head><title>" + ovr.race.Name + " Overall Results</title></head><body><pre>\n")
	if err != nil {
		return err
	}

	w := io.MultiWriter(f, os.Stdout)
	fmt.Printf("%s", "\x1Bc") // clear stdout
	fmt.Printf("Last Updated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "\n\n\n")
	fmt.Fprintln(w, "Place Bib   Name                             Grade Team                             Time")
	fmt.Fprintln(w, "===== ===== ================================ ===== ================================ ========")
	for _, r := range overallResults {
		fmt.Fprintf(w, "%-5d %-5d %-32s %-5d %-32s %-8s\n", r.Place, r.Bib, r.Athlete.Name(), r.Athlete.Grade, r.Athlete.Team, formatFinishTime(r.Finishtime))
	}

	_, err = f.WriteString("\n</pre></body></html>\n")
	if err != nil {
		return err
	}
	ovr.logger.Info("Done Building overall")

	return nil
}
