package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Mode string

const (
	Old Mode = "Old"
	New Mode = "New"

	Avg string = "Average"
	P50 string = "Perc50"
	P90 string = "Perc90"
	P99 string = "Perc99"
)

// A Table is a table for display in the output.
type Table struct {
	Mode Mode
	Keys map[string]*Row
	Rows []*Row // Ordered by input.
}

func (t *Table) LoadFile(file string, mode Mode) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := t.AddBenchmarkResult(scanner.Bytes(), mode); err != nil {
			log.Fatal(err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func getKey(metric, name, q string) string {
	return fmt.Sprintf("%v:%v:%v", metric, name, q)
}

func (t *Table) AddBenchmarkResult(bytes []byte, mode Mode) error {
	var dataItems DataItems
	if err := json.Unmarshal(bytes, &dataItems); err != nil {
		return err
	}

	for _, item := range dataItems.DataItems {
		for _, q := range []string{Avg, P50, P90, P99} {
			key := getKey(item.Labels["Metric"], item.Labels["Name"], q)
			var row *Row
			if r, ok := t.Keys[key]; ok {
				row = r
			} else {
				row = &Row{}
				t.Keys[key] = row
				t.Rows = append(t.Rows, row)

				row.Metric, row.Group = item.Labels["Metric"], item.Labels["Name"]
				row.Unit = item.Unit
				row.Quantizer = q
			}
			row.Add(item.Data, q, mode)
		}
	}

	return nil
}

// A Row is a table row for display in the output.
type Row struct {
	Metric    string   // benchmark name
	Group     string   // group name
	Quantizer string   // q-th quantile, or average
	Unit      string   // pods/s, or ms
	Old       *Metrics // old data
	New       *Metrics // new data
	Delta     string   // percent change
}

func (r *Row) Add(data map[string]float64, q string, mode Mode) {
	var metrics *Metrics
	if mode == Old {
		if r.Old == nil {
			r.Old = &Metrics{}
		}
		metrics = r.Old
	} else if mode == New {
		if r.New == nil {
			r.New = &Metrics{}
		}
		metrics = r.New
	}

	metrics.Values = append(metrics.Values, data[q])
}
