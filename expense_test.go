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
			Month:       "January",
			Year:        "2018",
		},
		"Description: Vacation in the carribean | Amount: 2000 | Category: optional | Date: January 2018",
	},
	{
		k.Expense{
			ID:          100,
			Description: "Pasta",
			Category:    "survival",
			Amount:      10,
			Month:       "October",
			Year:        "2017",
		},
		"Description: Pasta | Amount: 10 | Category: survival | Date: October 2017",
	},
	{
		k.Expense{
			ID:          2021023,
			Description: "Nike AirMax '97",
			Category:    "extra",
			Amount:      178,
			Month:       "September",
			Year:        "2017",
		},
		"Description: Nike AirMax '97 | Amount: 178 | Category: extra | Date: September 2017",
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
