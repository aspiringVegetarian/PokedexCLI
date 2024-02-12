package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HeLlo WoRlD",
			expected: []string{"hello", "world"},
		},
		{
			input:    "    hello     world     ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, cs := range cases {
		actual := cleanInput(cs.input)
		if len(actual) != len(cs.expected) {
			t.Errorf("The lengths are not equal: %v vs %v",
				len(actual),
				len(cs.expected))
			continue
		}

		for i := range actual {
			if actual[i] != cs.expected[i] {
				t.Errorf("The words do not match: %v vs %v",
					actual[i],
					cs.expected[i])
			}
		}

	}
}
