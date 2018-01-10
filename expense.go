package main

import "strconv"

type Expense struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Amount      int    `json:"amount"`
	Week        int    `json:"week"`
}

func (e Expense) String() string {
	buffer := " | "

	out := ""
	out += "Description: " + e.Description
	out += buffer
	out += "Amount: " + strconv.Itoa(e.Amount)
	out += buffer
	out += "Category: " + e.Category
	out += buffer
	out += "Week: " + strconv.Itoa(e.Week)

	return out
}
