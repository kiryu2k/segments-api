package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseResponse(t *testing.T) {
	type testCase struct {
		input       []byte
		expectedLhs uint64
		expectedRhs string
	}
	testCases := []testCase{
		{
			input:       []byte(""),
			expectedLhs: 0,
			expectedRhs: "",
		},
		{
			input:       []byte("()"),
			expectedLhs: 0,
			expectedRhs: "",
		},
		{
			input:       []byte("(8921,tesT)"),
			expectedLhs: 8921,
			expectedRhs: "tesT",
		},
		{
			input:       []byte("(321312, sdadads, 231123)"),
			expectedLhs: 0,
			expectedRhs: "",
		},
	}
	for _, test := range testCases {
		lhs, rhs := ParseResponse(test.input)
		assert.Equal(t, test.expectedLhs, lhs)
		assert.Equal(t, test.expectedRhs, rhs)
	}
}

func Test_ParseTTL(t *testing.T) {
	type testCase struct {
		input    string
		expected *TTL
	}
	testCases := []testCase{
		{
			input: "1y2m9d",
			expected: &TTL{
				Years:  1,
				Months: 2,
				Days:   9,
			},
		},
		{
			input:    "-1y-2m-9d",
			expected: nil,
		},
		{
			input: "100y11m30d",
			expected: &TTL{
				Years:  100,
				Months: 11,
				Days:   30,
			},
		},
		{
			input:    "100y11m31d",
			expected: nil,
		},
		{
			input:    "1y12m0d",
			expected: nil,
		},
		{
			input:    "101y0m31d",
			expected: nil,
		},
		{
			input: "64y",
			expected: &TTL{
				Years:  64,
				Months: 0,
				Days:   0,
			},
		},
		{
			input: "8m24d",
			expected: &TTL{
				Years:  0,
				Months: 8,
				Days:   24,
			},
		},
	}
	for _, test := range testCases {
		ttl, _ := ParseTTL(test.input)
		assert.Equal(t, test.expected, ttl)
	}
}
