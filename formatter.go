package ghaprofiler

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
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
		table := tablewriter.NewWriter(w)
		table.SetAutoFormatHeaders(false)
		table.SetHeader([]string{"Number", "Min", "Median", "Mean", "P50", "P90", "P95", "P99", "Max", "Name"})
		for _, p := range p.Profile {
			table.Append([]string{
				strconv.FormatInt(p.Number, 10),
				strconv.FormatFloat(p.Min, 'f', 6, 64),
				strconv.FormatFloat(p.Median, 'f', 6, 64),
				strconv.FormatFloat(p.Mean, 'f', 6, 64),
				strconv.FormatFloat(p.Percentile50, 'f', 6, 64),
				strconv.FormatFloat(p.Percentile90, 'f', 6, 64),
				strconv.FormatFloat(p.Percentile95, 'f', 6, 64),
				strconv.FormatFloat(p.Percentile99, 'f', 6, 64),
				strconv.FormatFloat(p.Max, 'f', 6, 64),
				p.Name,
			})
		}
		fmt.Fprintf(w, "Job: %s\n", p.Name)
		table.Render()
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
