package gonx

import "strconv"

// Reducer interface for Entries channel redure.
//
// Each Reduce method should accept input channel of Entries, do it's job and
// the result should be written to the output channel.
//
// It does not return values because usually it runs in a separate
// goroutine and it is handy to use channel for reduced data retrieval.
type Reducer interface {
	Reduce(input chan Entry, output chan interface{})
}

// Implements Reducer interface for simple input entries redirection to
// the output channel.
type ReadAll struct {
}

// Redirect input Entries channel directly to the output without any
// modifications. It is useful when you want jast to read file fast
// using asynchronous with mapper routines.
func (r *ReadAll) Reduce(input chan Entry, output chan interface{}) {
	output <- input
}

// Implements Reducer interface to count entries
type Count struct {
}

// Simply count entrries and write a sum to the output channel
func (r *Count) Reduce(input chan Entry, output chan interface{}) {
	count := 0
	for {
		_, ok := <-input
		if !ok {
			break
		}
		count++
	}
	output <- count
}

// Implements Reducer interface for summarize Entry values for the given fields
type Sum struct {
	Fields []string
}

// Summarize given Entry fields and return a map with result for each field.
func (r *Sum) Reduce(input chan Entry, output chan interface{}) {
	sum := make(map[string]float64)
	for _, name := range r.Fields {
		sum[name] = 0
	}
	for entry := range input {
		for _, name := range r.Fields {
			val, err := strconv.ParseFloat(entry[name], 64)
			if err == nil {
				sum[name] += val
			}
		}
	}
	output <- sum
}
