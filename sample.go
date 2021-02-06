package main

// Copied and tweaked from golang.org/x/perf/internal.

import (
	"math"
	"sort"
)

// Sample is a collection of data points.
type Sample struct {
	// Xs is the slice of sample values.
	Xs []float64

	// Sorted indicates that Xs is sorted in ascending order.
	Sorted bool
}

// Bounds returns the minimum and maximum values of the Sample.
func (s Sample) Bounds() (min, max float64) {
	if len(s.Xs) == 0 || !s.Sorted {
		return Bounds(s.Xs)
	}

	return s.Xs[0], s.Xs[len(s.Xs)-1]
}

// Bounds returns the minimum and maximum values of xs.
func Bounds(xs []float64) (min, max float64) {
	if len(xs) == 0 {
		return math.NaN(), math.NaN()
	}
	min, max = xs[0], xs[0]
	for _, x := range xs {
		if x < min {
			min = x
		}
		if x > max {
			max = x
		}
	}
	return
}

// Mean returns the arithmetic mean of xs.
func Mean(xs []float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	}
	m := 0.0
	for i, x := range xs {
		m += (x - m) / float64(i+1)
	}
	return m
}

// Percentile returns the pctileth value from the Sample. This uses
// interpolation method R8 from Hyndman and Fan (1996).
//
// pctile will be capped to the range [0, 1]. If len(xs) == 0, returns NaN.
//
// Percentile(0.5) is the median. Percentile(0.25) and
// Percentile(0.75) are the first and third quartiles, respectively.
func (s Sample) Percentile(pctile float64) float64 {
	if len(s.Xs) == 0 {
		return math.NaN()
	} else if pctile <= 0 {
		min, _ := s.Bounds()
		return min
	} else if pctile >= 1 {
		_, max := s.Bounds()
		return max
	}

	if !s.Sorted {
		s = *s.Copy().Sort()
	}

	N := float64(len(s.Xs))
	n := 1/3.0 + pctile*(N+1/3.0) // R8
	kf, frac := math.Modf(n)
	k := int(kf)
	if k <= 0 {
		return s.Xs[0]
	} else if k >= len(s.Xs) {
		return s.Xs[len(s.Xs)-1]
	}
	return s.Xs[k-1] + frac*(s.Xs[k]-s.Xs[k-1])
}

// IQR returns the interquartile range of the Sample.
//
// This is constant time if s.Sorted and s.Weights == nil.
func (s Sample) IQR() float64 {
	if !s.Sorted {
		s = *s.Copy().Sort()
	}
	return s.Percentile(0.75) - s.Percentile(0.25)
}

// Copy returns a copy of the Sample.
//
// The returned Sample shares no data with the original, so they can
// be modified (for example, sorted) independently.
func (s Sample) Copy() *Sample {
	xs := make([]float64, len(s.Xs))
	copy(xs, s.Xs)
	return &Sample{xs, s.Sorted}
}

// Sort sorts the samples in place in s and returns s.
//
// A sorted sample improves the performance of some algorithms.
func (s *Sample) Sort() *Sample {
	if s.Sorted || sort.Float64sAreSorted(s.Xs) {
		// All set
	} else {
		sort.Float64s(s.Xs)
	}
	s.Sorted = true
	return s
}
