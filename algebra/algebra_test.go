package algebra

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/markwiat/range-subtractor/span"
)

type finalResult int

func (ft finalResult) Add(other span.SpanLength) span.SpanLength {
	return ft + other.(finalResult)
}

type testCorner int64

func (tc testCorner) Sub(other span.Corner) span.SpanLength {
	return finalResult(tc - other.(testCorner))
}

func (tc testCorner) Before(other span.Corner) bool {
	return tc < other.(testCorner)
}

type testSpan struct {
	start testCorner
	end   testCorner
}

func (ts testSpan) Start() span.Corner {
	return ts.start
}

func (ts testSpan) End() span.Corner {
	return ts.end
}

type categorizedTestSpan struct {
	span.Span
	super bool
}

func (cts categorizedTestSpan) IsSuper() bool {
	return cts.super
}

func newCategorizedSpan(start, end int, super bool) span.CategorizedSpan {
	return categorizedTestSpan{
		Span:  testSpan{start: testCorner(start), end: testCorner(end)},
		super: super,
	}
}

func mapToTestSpans(spans []span.Span) []span.Span {
	if spans == nil {
		return []span.Span{}
	}
	result := make([]span.Span, len(spans))
	for i, span := range spans {
		result[i] = testSpan{
			start: span.Start().(testCorner),
			end:   span.End().(testCorner),
		}
	}
	return result
}

func TestFiltering(t *testing.T) {
	positive1 := testSpan{1, 4}
	positive2 := testSpan{-2, -1}
	negative := testSpan{2, 1}
	null := testSpan{5, 5}
	spans := []span.Span{positive1, negative, positive2, null}

	filtered := FilterOutNotPositive(spans)

	expected := []span.Span{positive1, positive2}
	assert.Equal(t, filtered, expected)
}

func TestSorting(t *testing.T) {
	a := testSpan{1, 4}
	b1 := testSpan{2, 4}
	b2 := testSpan{2, 5}
	c := testSpan{3, 6}
	d := testSpan{7, 8}

	expected1 := []span.Span{a, b1, b2, c, d}
	expected2 := []span.Span{a, b2, b1, c, d}

	isExpected := func(result []span.Span) bool {
		return assert.IsEqual(result, expected1) || assert.IsEqual(result, expected2)
	}

	inputs := [][]span.Span{
		{d, c, b2, b1, a},
		{a, c, b2, b1, d},
		{a, b2, b1, c, d},
		{c, b2, d, a, b1},
	}

	for _, input := range inputs {
		result := SortByStart(input)
		if !isExpected(result) {
			assert.Equal(t, result, true)
		}
	}
}

func TestJoin(t *testing.T) {
	a := testSpan{-5, -3}
	b := testSpan{-1, 2}
	c := testSpan{0, 3}
	d := testSpan{3, 7}
	e := testSpan{6, 9}
	f := testSpan{11, 12}

	expectedMerged := testSpan{-1, 9}

	inputs := [][]span.Span{
		{f, e, d, a, b, c},
		{c, c, e, d, b, c, b},
	}
	expected := [][]span.Span{
		{a, expectedMerged, f},
		{expectedMerged},
	}

	for i, input := range inputs {
		result := JoinOverlapped(input)
		result = mapToTestSpans(result)
		exp := expected[i]
		assert.Equal(t, result, exp)
	}
}

func TestSubtractSpans(t *testing.T) {
	type testData struct {
		super      []span.Span
		subtrahend []span.Span
		expected   []span.Span
	}

	tests := []testData{
		{
			super:      []span.Span{testSpan{-4, -3}, testSpan{0, 20}, testSpan{21, 25}, testSpan{26, 30}, testSpan{31, 35}},
			subtrahend: []span.Span{testSpan{-1, 2}, testSpan{4, 8}, testSpan{10, 12}, testSpan{19, 22}},
			expected:   []span.Span{testSpan{-4, -3}, testSpan{2, 4}, testSpan{8, 10}, testSpan{12, 19}, testSpan{22, 25}, testSpan{26, 30}, testSpan{31, 35}},
		},
		{
			super:      []span.Span{testSpan{2, 10}, testSpan{20, 100}},
			subtrahend: []span.Span{testSpan{0, 1}, testSpan{30, 40}, testSpan{50, 60}},
			expected:   []span.Span{testSpan{2, 10}, testSpan{20, 30}, testSpan{40, 50}, testSpan{60, 100}},
		},
		{
			super:      []span.Span{testSpan{2, 10}, testSpan{20, 100}},
			subtrahend: []span.Span{testSpan{0, 1}, testSpan{11, 14}, testSpan{30, 40}, testSpan{50, 60}},
			expected:   []span.Span{testSpan{2, 10}, testSpan{20, 30}, testSpan{40, 50}, testSpan{60, 100}},
		},
		{
			super:      []span.Span{testSpan{2, 10}, testSpan{20, 100}, testSpan{300, 400}},
			subtrahend: []span.Span{testSpan{0, 1}, testSpan{11, 14}, testSpan{15, 18}, testSpan{30, 40}, testSpan{50, 60}, testSpan{150, 160}, testSpan{170, 180}, testSpan{280, 320}},
			expected:   []span.Span{testSpan{2, 10}, testSpan{20, 30}, testSpan{40, 50}, testSpan{60, 100}, testSpan{320, 400}},
		},
		{
			super:      []span.Span{testSpan{2, 10}, testSpan{20, 100}},
			subtrahend: []span.Span{testSpan{0, 1}, testSpan{-5, -2}},
			expected:   []span.Span{testSpan{2, 10}, testSpan{20, 100}},
		},
		{
			super:      []span.Span{testSpan{2, 10}, testSpan{20, 100}},
			subtrahend: []span.Span{testSpan{121, 125}, testSpan{150, 160}},
			expected:   []span.Span{testSpan{2, 10}, testSpan{20, 100}},
		},
		{
			super:      []span.Span{testSpan{-4, -3}, testSpan{0, 3}, testSpan{5, 8}, testSpan{10, 15}},
			subtrahend: []span.Span{testSpan{-2, 14}},
			expected:   []span.Span{testSpan{-4, -3}, testSpan{14, 15}},
		},
		{
			super:      []span.Span{testSpan{-4, -3}, testSpan{0, 3}, testSpan{5, 8}, testSpan{10, 15}},
			subtrahend: []span.Span{testSpan{-4, 15}},
			expected:   []span.Span{},
		},
		{
			super:      []span.Span{testSpan{-4, -3}, testSpan{0, 3}, testSpan{5, 8}, testSpan{10, 15}},
			subtrahend: []span.Span{testSpan{-5, 16}},
			expected:   []span.Span{},
		},
	}

	for _, test := range tests {
		spans := make([]span.CategorizedSpan, 0, len(test.super)+len(test.super))
		for _, s := range test.super {
			spans = append(spans, categorizedTestSpan{s, true})
		}
		for _, s := range test.subtrahend {
			spans = append(spans, categorizedTestSpan{s, false})
		}
		result := FindSubtractedSpans(spans)
		result = mapToTestSpans(result)

		assert.Equal(t, result, test.expected)
	}
}

func TestSubtractLength(t *testing.T) {
	type testData struct {
		spans    []span.CategorizedSpan
		expected int
	}

	emptyLength := finalResult(0)
	tests := []testData{
		{
			spans: []span.CategorizedSpan{newCategorizedSpan(0, 10, true), newCategorizedSpan(3, 5, false),
				newCategorizedSpan(3, 5, false), newCategorizedSpan(4, 11, true)},
			expected: 9,
		},
		{
			spans: []span.CategorizedSpan{newCategorizedSpan(0, 10, true), newCategorizedSpan(3, 5, false),
				newCategorizedSpan(-1, 11, false), newCategorizedSpan(4, 11, true)},
			expected: 0,
		},
		{
			spans: []span.CategorizedSpan{newCategorizedSpan(-5, 3, true), newCategorizedSpan(-2, 3, true),
				newCategorizedSpan(-6, 4, true), newCategorizedSpan(3, 4, false), newCategorizedSpan(4, 8, false),
				newCategorizedSpan(7, 11, false), newCategorizedSpan(10, 12, true), newCategorizedSpan(14, 15, true)},
			expected: 11,
		},
		{
			spans: []span.CategorizedSpan{newCategorizedSpan(-2, 3, false), newCategorizedSpan(10, 15, false),
				newCategorizedSpan(0, 30, true), newCategorizedSpan(4, 6, false), newCategorizedSpan(27, 30, false)},
			expected: 17,
		},
		{
			spans:    []span.CategorizedSpan{newCategorizedSpan(-2, 3, false), newCategorizedSpan(10, 15, false)},
			expected: 0,
		},
		{
			spans:    []span.CategorizedSpan{newCategorizedSpan(2, 2, true)},
			expected: 0,
		},
		{
			spans:    []span.CategorizedSpan{},
			expected: 0,
		},
		{
			spans:    []span.CategorizedSpan{newCategorizedSpan(2, 3, true), newCategorizedSpan(-5, -1, true)},
			expected: 5,
		},
		{
			spans: []span.CategorizedSpan{newCategorizedSpan(-2, 6, true), newCategorizedSpan(5, 6, false),
				newCategorizedSpan(9, 10, false), newCategorizedSpan(12, 14, false)},
			expected: 7,
		},
		{
			spans: []span.CategorizedSpan{newCategorizedSpan(-2, 6, true), newCategorizedSpan(-6, -2, false),
				newCategorizedSpan(6, 10, false), newCategorizedSpan(10, 12, true)},
			expected: 10,
		},
	}

	for _, test := range tests {
		result := SubtractFromSuperSpans(emptyLength, test.spans)
		r := result.(finalResult)
		rInt := int(r)
		assert.Equal(t, rInt, test.expected)
	}
}
