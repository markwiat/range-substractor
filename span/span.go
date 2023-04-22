package span

type SpanLength interface {
	Sum(other SpanLength) SpanLength
}

type Corner interface {
	Before(other Corner) bool
	Sub(other Corner) SpanLength
}

type Span interface {
	Start() Corner
	End() Corner
}

type CategorizedSpan interface {
	Span
	IsSuper() bool
}

func FilterOutNotPositive(spans []Span) []Span {
	return nil
}

func SortByStart(spans []Span) []Span {
	return nil
}

func JoinOverlapped(spans []Span) []Span {
	return nil
}

func SubtractFromSuperSpans(spans []CategorizedSpan) SpanLength {
	return nil
}
