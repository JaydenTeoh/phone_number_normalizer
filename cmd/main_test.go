package main

import "testing"

func TestNormalize(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"123-456-7890", "1234567890"},
		{"1234567892", "1234567892"},
		{"(123)456-7892", "1234567892"},
	}
	for _, testcase := range testCases {
		t.Run(testcase.input, func(t *testing.T) {
			actual := normalize(testcase.input)
			if actual != testcase.want {
				t.Errorf("got %s; want %s", actual, testcase.want)
			}
		})
	}
}
