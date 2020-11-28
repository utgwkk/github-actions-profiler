package ghaprofiler

import (
	"github.com/google/go-github/v32/github"
	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
)

type TaskStepProfile struct {
	Name         string  `json:"name"`
	Number       int64   `json:"number"`
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	Median       float64 `json:"median"`
	Mean         float64 `json:"mean"`
	Percentile95 float64 `json:"95percentile"`
	Percentile99 float64 `json:"99percentile"`
}

type TaskStepProfileResult = []*TaskStepProfile

func ProfileTaskStep(steps []*github.TaskStep) (profileResult TaskStepProfileResult, err error) {
	taskStepsByNumber := make(map[int64][]*github.TaskStep)

	// aggregate tasks by its number
	for _, step := range steps {
		taskStepsByNumber[*step.Number] = append(taskStepsByNumber[*step.Number], step)
	}

	for stepNumber, steps := range taskStepsByNumber {
		var stepElapsed []float64
		stepName := steps[0].Name

		for _, step := range steps {
			elapsed := step.CompletedAt.Sub(step.StartedAt.Time)
			elapsedSeconds := float64(elapsed.Nanoseconds()) / 1e9
			stepElapsed = append(stepElapsed, elapsedSeconds)
		}

		min, err := stats.Min(stepElapsed)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate min")
		}
		max, err := stats.Max(stepElapsed)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate max")
		}
		median, err := stats.Median(stepElapsed)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate median")
		}
		mean, err := stats.Mean(stepElapsed)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate mean")
		}
		percentile99, err := stats.Percentile(stepElapsed, 99)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate 99%ile")
		}
		percentile95, err := stats.Percentile(stepElapsed, 95)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate 95%ile")
		}

		profileResult = append(profileResult, &TaskStepProfile{
			Name:         *stepName,
			Number:       stepNumber,
			Min:          min,
			Max:          max,
			Median:       median,
			Mean:         mean,
			Percentile95: percentile95,
			Percentile99: percentile99,
		})
	}

	return
}
