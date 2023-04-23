package algebra

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/markwiat/range-substractor/span"
)

type finalResult int

func (ft finalResult) Sum(other span.SpanLength) span.SpanLength {
	r := ft + other.(finalResult)
	return r
}

type testCorner int64

func (tc testCorner) Sub(other span.Corner) span.SpanLength {
	otc := other.(testCorner)
	return finalResult(tc - otc)
}

func (tc testCorner) Before(other span.Corner) bool {
	otc := other.(testCorner)
	return tc < otc
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

func mapToTestSpans(spans []span.Span) []span.Span {
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
	b := testSpan{2, 4}
	b2 := testSpan{2, 5}
	c := testSpan{3, 6}
	d := testSpan{7, 8}

	expected1 := []span.Span{a, b, b2, c, d}
	expected2 := []span.Span{a, b2, b, c, d}

	isExpected := func(result []span.Span) bool {
		return assert.IsEqual(result, expected1) || assert.IsEqual(result, expected2)
	}

	inputs := [][]span.Span{
		{d, c, b2, b, a},
		{a, c, b2, b, d},
		{a, b2, b, c, d},
		{c, b2, d, a, b},
	}

	for i, input := range inputs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result := SortByStart(input)
			if !isExpected(result) {
				assert.Equal(t, result, true)
			}
		})
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

func TestSubtractOrdered(t *testing.T) {
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
	}

	for _, test := range tests {
		spans := make([]span.CategorizedSpan, 0, len(test.super)+len(test.super))
		for _, s := range test.super {
			spans = append(spans, categorizedTestSpan{s, true})
		}
		for _, s := range test.subtrahend {
			spans = append(spans, categorizedTestSpan{s, false})
		}
		result := SubtractAndGet(spans)
		result = mapToTestSpans(result)

		assert.Equal(t, result, test.expected)
	}
}
