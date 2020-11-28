package ghaprofiler

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const (
	formatNameJSON  = "json"
	formatNameTable = "table"
)

var availableFormats = []string{
	formatNameJSON,
	formatNameTable,
}

func AvailableFormatsForCLI() string {
	return strings.Join(availableFormats, ", ")
}

type ProfileForFormatter struct {
	Name    string             `json:"name"`
	Profile []*TaskStepProfile `json:"profile"`
}

type ProfileInput []*ProfileForFormatter

func IsValidFormatName(formatName string) bool {
	for _, available := range availableFormats {
		if formatName == available {
			return true
		}
	}
	return false
}

func WriteJSON(w io.Writer, profileResult ProfileInput) (err error) {
	encoder := json.NewEncoder(w)
	err = encoder.Encode(struct {
		Profiles []*ProfileForFormatter `json:"profiles"`
	}{
		Profiles: profileResult,
	})
	return
}

func WriteTable(w io.Writer, profileResult ProfileInput) error {
	for _, p := range profileResult {
		fmt.Fprintf(w, "Job: %s\n", p.Name)
		fmt.Fprintln(w, "Number\tMin\tMedian\tMean\tP95\tP99\tMax\tName")
		for _, p := range p.Profile {
			fmt.Fprintf(w, "%d\t%f\t%f\t%f\t%f\t%f\t%f\t%s\n", p.Number, p.Min, p.Median, p.Mean, p.Percentile95, p.Percentile99, p.Max, p.Name)
		}
		fmt.Fprintln(w)
	}
	return nil
}

func WriteWithFormat(w io.Writer, profileResult ProfileInput, format string) error {
	switch format {
	case formatNameJSON:
		WriteJSON(w, profileResult)
		break
	case formatNameTable:
		WriteTable(w, profileResult)
		break
	default:
		return fmt.Errorf("Invalid format: %s", format)
	}
	return nil
}
