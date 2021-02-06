package main

// A Metrics holds the measurements of a single data item
// for all runs of a particular benchmark.
type Metrics struct {
	// Unit    string    // unit being measured
	Values  []float64 // measured values
	RValues []float64 // Values with outliers removed
	Min     float64   // min of RValues
	Mean    float64   // mean of RValues
	Max     float64   // max of RValues
}

// compute updates the derived statistics in m from the raw
// samples in m.Values.
func (m *Metrics) compute() {
	// Discard outliers.
	values := Sample{Xs: m.Values}
	q1, q3 := values.Percentile(0.25), values.Percentile(0.75)
	lo, hi := q1-1.5*(q3-q1), q3+1.5*(q3-q1)
	for _, value := range m.Values {
		if lo <= value && value <= hi {
			m.RValues = append(m.RValues, value)
		}
	}

	// Compute statistics of remaining data.
	m.Min, m.Max = Bounds(m.RValues)
	m.Mean = Mean(m.RValues)
}
