package span

type SpanLength interface {
	Add(other SpanLength) SpanLength
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
