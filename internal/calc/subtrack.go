package calc

import "github.com/markwiat/range-substractor/span"

func overlaps(super, subtrahend span.Span) bool {
	return subtrahend.Start().Before(super.End())
}

type spanPosition int

const (
	left spanPosition = iota
	leftEdge
	inside
	rightEdge
	right
)

type subtractResultType int

const (
	empty subtractResultType = iota
	whole
	leftSide
	rightSide
	bothSide
)

func SubtractOrdered(supers, subtrahends []span.Span) []span.Span {
	result := make([]span.Span, 0)

	i, j := 0, 0
	for ; i < len(supers) && j < len(subtrahends); i++ {
		subResult := subtractWithAll(supers[i], subtrahends, &j)
		result = append(result, subResult...)
	}

	for ; i < len(supers); i++ {
		result = append(result, supers[i])
	}

	return result
}

func SplitByCategory(spans []span.CategorizedSpan) (supers, subtrahends []span.Span) {
	for _, span := range spans {
		if span.IsSuper() {
			supers = append(supers, span)
		} else {
			subtrahends = append(subtrahends, span)
		}
	}

	return
}

func SumLengths(spans []span.Span) span.SpanLength {
	if len(spans) == 0 {
		return nil
	}
	result := spanLength(spans[0])
	for i := 1; i < len(spans); i++ {
		result = result.Sum(spanLength(spans[i]))
	}

	return result
}

func spanLength(span span.Span) span.SpanLength {
	return span.End().Sub(span.Start())
}

func subtractWithAll(super span.Span, subtrahends []span.Span, index *int) []span.Span {
	result := make([]span.Span, 0)
	for ; *index < len(subtrahends); *index++ {
		l, r := subtractOne(super, subtrahends[*index])
		if l != nil {
			result = append(result, l)
		}
		if r != nil {
			*index++
			subResult := subtractWithAll(r, subtrahends, index)
			if len(subResult) == 0 {
				result = append(result, r)
			} else {
				result = append(result, subResult...)
			}
			break
		}
		if l != nil || !subtrahends[*index].End().Before(super.End()) {
			break
		}
	}

	return result
}

func subtractOne(super, subtrahend span.Span) (left, right span.Span) {
	setLeft := func() {
		left = &spanImpl{
			start: super.Start(),
			end:   subtrahend.Start(),
		}
	}
	setRight := func() {
		right = &spanImpl{
			start: subtrahend.End(),
			end:   super.End(),
		}
	}

	subtractType := findSubtractType(super, subtrahend)

	switch subtractType {
	case empty:
		return
	case whole:
		left = super
	case leftSide:
		setLeft()
	case rightSide:
		setRight()
	case bothSide:
		setLeft()
		setRight()
	}
	return
}

func findSubtractType(super, subtrahend span.Span) subtractResultType {
	startPos := findPosition(super, subtrahend.Start())
	endPos := findPosition(super, subtrahend.End())

	if startPos == rightEdge || startPos == right || endPos == left || endPos == leftEdge {
		return whole
	} else if startPos == inside && endPos == inside {
		return bothSide
	} else if startPos == inside {
		return leftSide
	} else if endPos == inside {
		return rightSide
	}
	return empty
}

func findPosition(span span.Span, corner span.Corner) spanPosition {
	if corner.Before(span.Start()) {
		return left
	} else if equals(corner, span.Start()) {
		return leftEdge
	} else if equals(corner, span.End()) {
		return rightEdge
	} else if span.End().Before(corner) {
		return right
	}
	return inside
}
