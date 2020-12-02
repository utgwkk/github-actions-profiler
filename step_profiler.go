package ghaprofiler

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
)

type TaskStepProfile struct {
	Name        string                    `json:"name"`
	Number      int64                     `json:"number"`
	Min         float64                   `json:"min"`
	Max         float64                   `json:"max"`
	Median      float64                   `json:"median"`
	Mean        float64                   `json:"mean"`
	Percentiles map[int64]*percentileData `json:"percentiles"`
}

var percentiles = []int64{
	50,
	90,
	95,
	99,
}

type percentileData struct {
	Percentile int64
	Value      float64
}

func (pd *percentileData) MarshalJSON() ([]byte, error) {
	return json.Marshal(pd.Value)
}

func (pd percentileData) Label() string {
	return fmt.Sprintf("P%d", pd.Percentile)
}

func (pd percentileData) SortKey() string {
	return fmt.Sprintf("p%d", pd.Percentile)
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
		percentileResult := map[int64]*percentileData{}
		for _, percentile := range percentiles {
			value, err := stats.Percentile(stepElapsed, float64(percentile))
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed to calculate %d%%ile", percentile))
			}
			percentileResult[percentile] = &percentileData{Percentile: percentile, Value: value}
		}

		profileResult = append(profileResult, &TaskStepProfile{
			Name:        *stepName,
			Number:      stepNumber,
			Min:         min,
			Max:         max,
			Median:      median,
			Mean:        mean,
			Percentiles: percentileResult,
		})
	}

	return
}
