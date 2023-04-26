# range-subtractor
Find difference between ranges, even overlapped, e.g times or distances

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/markwiat/range-subtractor/algebra"
	"github.com/markwiat/range-subtractor/span"
)

type durationAlias time.Duration

func (d durationAlias) Add(other span.SpanLength) span.SpanLength {
	return d + other.(durationAlias)
}

type timeAlias time.Time

func (t timeAlias) Sub(other span.Corner) span.SpanLength {
	return durationAlias(time.Time(t).Sub(time.Time(other.(timeAlias))))
}

func (t timeAlias) Before(other span.Corner) bool {
	return time.Time(t).Before(time.Time(other.(timeAlias)))
}

type TimeRange struct {
	start timeAlias
	end   timeAlias
}

func (tr TimeRange) Start() span.Corner {
	return tr.start
}

func (tr TimeRange) End() span.Corner {
	return tr.end
}

type CategorizedTimeRange struct {
	TimeRange
	super bool
}

func (ctr CategorizedTimeRange) IsSuper() bool {
	return ctr.super
}

func createCategorizedTimeRange(super bool, start, end string) span.CategorizedSpan {
	ts, err := time.Parse(time.RFC3339, start)
	if err != nil {
		panic(err)
	}
	te, err := time.Parse(time.RFC3339, end)
	if err != nil {
		panic(err)
	}
	var ctr CategorizedTimeRange
	ctr.start = timeAlias(ts)
	ctr.end = timeAlias(te)
	ctr.super = super

	return ctr
}

func printRange(span span.Span) {
	if span == nil {
		fmt.Printf("No result\n")
	}
	start := time.Time(span.Start().(timeAlias))
	end := time.Time(span.End().(timeAlias))

	fmt.Printf("start: %v, end: %v\n", start, end)
}

func main() {
	super1 := createCategorizedTimeRange(true, "2023-05-01T08:00:00Z", "2023-05-01T12:00:00Z")
	super2 := createCategorizedTimeRange(true, "2023-05-01T10:00:00Z", "2023-05-01T14:00:00Z")
	subtrahend1 := createCategorizedTimeRange(false, "2023-05-01T09:30:00Z", "2023-05-01T10:15:00Z")
	subtrahend2 := createCategorizedTimeRange(false, "2023-05-01T10:00:00Z", "2023-05-01T10:30:00Z")

	spans := []span.CategorizedSpan{super1, super2, subtrahend1, subtrahend2}

	subtractedSpans := algebra.FindSubtractedSpans(spans)

	fmt.Printf("Subtracted spans:\n")
	for _, s := range subtractedSpans {
		printRange(s)
	}

	var emptyLength durationAlias
	finalLength := algebra.SubtractFromSuperSpans(emptyLength, spans)
	duration := time.Duration(finalLength.(durationAlias))

	fmt.Printf("Final result: %v\n", duration)
}
```
