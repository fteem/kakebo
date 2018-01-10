package main

import "strconv"

type Expense struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Amount      int    `json:"amount"`
	Month       string `json:"month"`
	Year        string `json:"year"`
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
	out += "Date: " + e.Month + " " + e.Year

	return out
}
