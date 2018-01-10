package main_test

import (
	k "github.com/fteem/kakebo"
	"testing"
)

var stringTests = []struct {
	input    k.Expense
	expected string
}{
	{
		k.Expense{
			ID:          1,
			Description: "Vacation in the carribean",
			Category:    "optional",
			Amount:      2000,
			Week:        25,
		},
		"Description: Vacation in the carribean | Amount: 2000 | Category: optional | Week: 25",
	},
	{
		k.Expense{
			ID:          100,
			Description: "Pasta",
			Category:    "survival",
			Amount:      10,
			Week:        2,
		},
		"Description: Pasta | Amount: 10 | Category: survival | Week: 2",
	},
	{
		k.Expense{
			ID:          2021023,
			Description: "Nike AirMax '97",
			Category:    "extra",
			Amount:      178,
			Week:        56,
		},
		"Description: Nike AirMax '97 | Amount: 178 | Category: extra | Week: 56",
	},
}

func TestString(t *testing.T) {
	for _, tt := range stringTests {
		actual := tt.input.String()
		if actual != tt.expected {
			t.Errorf("String(): expected '%s', actual '%s'", tt.expected, actual)
		}
	}
}
