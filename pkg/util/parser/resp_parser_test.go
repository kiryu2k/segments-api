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
