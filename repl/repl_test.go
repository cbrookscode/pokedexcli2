package repl

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  Hello WOrld    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  brown now COWS  ",
			expected: []string{"brown", "now", "cows"},
		},
		{
			input:    "  brown,    now! COWS. ",
			expected: []string{"brown", "now", "cows"},
		},
		{
			input:    "  brown,    now!COWS. ",
			expected: []string{"brown", "nowcows"},
		},
		{
			input:    "  b4rown3,    now!5COWS.1 ",
			expected: []string{"brown", "nowcows"},
		},
	}
	for i, c := range cases {
		actual := CleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Length of actual and expected slice dont match: Test: %v, Got: %v, Want: %v", i, actual, c.expected)
			t.Fail()
			continue
		}

		for index := range actual {
			if actual[index] != c.expected[index] {
				t.Errorf("actual word and expected word dont match: Test: %v, Got: %v, Want: %v", i, actual[index], c.expected[index])
				t.Fail()
				continue
			}
		}
	}
}
