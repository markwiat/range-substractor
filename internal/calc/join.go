package calc

import (
	"sort"

	"github.com/markwiat/range-substractor/span"
)

func FilterOutNotPositive(spans []span.Span) []span.Span {
	result := make([]span.Span, 0, len(spans))
	for _, span := range spans {
		if span.Start().Before(span.End()) {
			result = append(result, span)
		}
	}
	return result
}

func SortPositiveBySpanStart(spans []span.Span) []span.Span {
	result := make([]span.Span, len(spans))
	copy(result, spans)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Start().Before(result[j].Start())
	})
	return result
}

func JoinSorted(spans []span.Span) []span.Span {
	result := make([]span.Span, 0)
	for i := 0; i < len(spans); {
		joined, next := joinFirst(spans, i)
		result = append(result, joined)
		i = next
	}
	return result
}

type spanImpl struct {
	start span.Corner
	end   span.Corner
}

func (s *spanImpl) Start() span.Corner {
	return s.start
}

func (s *spanImpl) End() span.Corner {
	return s.end
}

func equals(a, b span.Corner) bool {
	return !a.Before(b) && !b.Before(a)
}

func aBeforeOrEqualsB(a, b span.Corner) bool {
	return !b.Before(a)
}

func aAfterB(a, b span.Corner) bool {
	return b.Before(a)
}

func aAfterOrEqualsB(a, b span.Corner) bool {
	return !a.Before(b)
}

func overlapsOrAdheresSorted(first, second span.Span) bool {
	return !first.End().Before(second.Start())
}

func joinFirst(spans []span.Span, index int) (span.Span, int) {
	result := spans[index]
	i := index + 1
	for ; i < len(spans); i++ {
		joined := joinOverlappedSorted(result, spans[i])
		if joined == nil {
			break
		}
		result = joined
	}
	return result, i
}

func joinOverlappedSorted(first, second span.Span) span.Span {
	if overlapsOrAdheresSorted(first, second) {
		return &spanImpl{
			start: first.Start(),
			end:   takeFurther(first.End(), second.End()),
		}
	}
	return nil
}

func takeFurther(a, b span.Corner) span.Corner {
	if a.Before(b) {
		return b
	}
	return a
}
