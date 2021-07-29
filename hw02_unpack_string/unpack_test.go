package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type tests struct {
	input    string
	expected string
}

func TestUnpack(t *testing.T) {
	tests := []tests{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestUnpackMyCases(t *testing.T) {
	tests := []tests{
		{input: `f\\\\\\\\\\\\\\1`, expected: `f\\\\\\\`},
		{input: `a2da\\\9`, expected: `aada\9`},
		{input: `a9b9c9\9\9ss`, expected: `aaaaaaaaabbbbbbbbbccccccccc99ss`},
		{input: `\\5\\\3\j\\\\`, expected: `\\\\\\3j\\`},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "d\n5\n\n\n\n\n\nnnnn\n6abc", expected: "d\n\n\n\n\n\n\n\n\n\n\nnnnn\n\n\n\n\n\nabc"},
	}
	runTests(t, tests)
}

func runTests(t *testing.T, tests []tests) {
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}
