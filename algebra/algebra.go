package algebra

import (
	"github.com/markwiat/range-substractor/internal/calc"
	"github.com/markwiat/range-substractor/span"
)

func FilterOutNotPositive(spans []span.Span) []span.Span {
	return calc.FilterOutNotPositive(spans)
}

func SortByStart(spans []span.Span) []span.Span {
	positives := FilterOutNotPositive(spans)
	return calc.SortPositiveBySpanStart(positives)
}

func JoinOverlapped(spans []span.Span) []span.Span {
	sorted := SortByStart(spans)
	return calc.JoinSorted(sorted)
}

func FindSubtractedSpans(spans []span.CategorizedSpan) []span.Span {
	supers, subtrahends := calc.SplitByCategory(spans)
	supers = JoinOverlapped(supers)
	subtrahends = JoinOverlapped(subtrahends)
	subtracted := calc.SubtractOrdered(supers, subtrahends)

	return subtracted
}

func SubtractFromSuperSpans(emptyResult span.SpanLength, spans []span.CategorizedSpan) span.SpanLength {
	subtracted := FindSubtractedSpans(spans)

	result := calc.SumLengths(subtracted)
	if result == nil {
		return emptyResult
	}

	return result
}
