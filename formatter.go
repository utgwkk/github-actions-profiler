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
	formatNameJSON     = "json"
	formatNameMarkdown = "markdown"
	formatNameTable    = "table"
	formatNameTSV      = "tsv"
)

var availableFormats = []string{
	formatNameJSON,
	formatNameMarkdown,
	formatNameTable,
	formatNameTSV,
}

var fields = []string{
	"number",
	"min",
	"median",
	"mean",
	"p50",
	"p90",
	"p95",
	"p99",
	"max",
	"name",
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

func showHeader(showIf filterFieldFunc) []string {
	var header []string
	for _, field := range fields {
		if showIf(field) {
			header = append(header, field)
		}
	}
	return header
}

func WriteTable(w io.Writer, profileResult ProfileInput, markdown bool, showIf filterFieldFunc) error {
	for _, p := range profileResult {
		table := tablewriter.NewWriter(w)
		table.SetAutoFormatHeaders(false)
		if markdown {
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
			table.SetAutoWrapText(false)
		}
		table.SetHeader(showHeader(showIf))
		for _, p := range p.Profile {
			var data []string
			if showIf("number") {
				data = append(data, strconv.FormatInt(p.Number, 10))
			}
			if showIf("min") {
				data = append(data, strconv.FormatFloat(p.Min, 'f', 6, 64))
			}
			if showIf("median") {
				data = append(data, strconv.FormatFloat(p.Median, 'f', 6, 64))
			}
			if showIf("mean") {
				data = append(data, strconv.FormatFloat(p.Mean, 'f', 6, 64))
			}
			if showIf("percentile50") {
				data = append(data, strconv.FormatFloat(p.Percentile50, 'f', 6, 64))
			}
			if showIf("percentile90") {
				data = append(data, strconv.FormatFloat(p.Percentile90, 'f', 6, 64))
			}
			if showIf("percentile95") {
				data = append(data, strconv.FormatFloat(p.Percentile95, 'f', 6, 64))
			}
			if showIf("percentile99") {
				data = append(data, strconv.FormatFloat(p.Percentile99, 'f', 6, 64))
			}
			if showIf("max") {
				data = append(data, strconv.FormatFloat(p.Max, 'f', 6, 64))
			}
			if showIf("name") {
				data = append(data, p.Name)
			}
			table.Append(data)
		}
		if markdown {
			fmt.Fprintf(w, "# Job: %s\n", p.Name)
			fmt.Fprintln(w)
		} else {
			fmt.Fprintf(w, "Job: %s\n", p.Name)
		}
		table.Render()
		fmt.Fprintln(w)
	}
	return nil
}

func WriteTSV(w io.Writer, profileResult ProfileInput) error {
	for _, p := range profileResult {
		fmt.Fprintf(w, "Job: %s\n", p.Name)
		fmt.Fprintln(w, "Number\tMin\tMedian\tMean\tP50\tP90\tP95\tP99\tMax\tName")
		for _, p := range p.Profile {
			fmt.Fprintf(w, "%d\t%f\t%f\t%f\t%f\t%f\t%f\t%f\t%f\t%s\n", p.Number, p.Min, p.Median, p.Mean, p.Percentile50, p.Percentile90, p.Percentile95, p.Percentile99, p.Max, p.Name)
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
		WriteTable(w, profileResult, false, func(_ string) bool { return true })
		break
	case formatNameMarkdown:
		WriteTable(w, profileResult, true, func(_ string) bool { return true })
		break
	case formatNameTSV:
		WriteTSV(w, profileResult)
		break
	default:
		return fmt.Errorf("Invalid format: %s", format)
	}
	return nil
}
