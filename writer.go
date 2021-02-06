package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var oneFileHeaders = []string{"Metric", "Group", "Quantile", "Unit", "Median"}
var twoFilesHeaders = []string{"Metric", "Group", "Quantile", "Unit", "Old", "New", "Diff"}

func FormatText(t *Table) {
	table := tablewriter.NewWriter(os.Stdout)
	if t.Mode == Old {
		table.SetHeader(oneFileHeaders)
	} else if t.Mode == New {
		table.SetHeader(twoFilesHeaders)
	}

	for _, row := range t.Rows {
		var data []string
		data = append(data, strings.TrimSuffix(row.Metric, "_duration_seconds"))
		// Shorten name field.
		sections := strings.Split(row.Group, "/")
		data = append(data, strings.Join(sections[1:len(sections)-1], "/"))
		data = append(data, row.Quantizer)
		data = append(data, row.Unit)
		// Discard outliers.
		row.Old.compute()
		data = append(data, fmt.Sprintf("%.2f", row.Old.Mean))
		if row.New != nil {
			row.New.compute()
			data = append(data, fmt.Sprintf("%.2f", row.New.Mean))
			x, y := row.Old.Mean, row.New.Mean
			diff := y/x - 1
			diffStr := fmt.Sprintf("%.2f%%", diff*100)
			if diff >= 0 {
				diffStr = fmt.Sprintf("+%s", diffStr)
			}
			data = append(data, diffStr)
		}

		table.Append(data)
	}

	table.Render()
}
