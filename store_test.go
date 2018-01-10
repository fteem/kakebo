package main_test

import (
	"testing"
	"time"

	"github.com/boltdb/bolt"
	k "github.com/fteem/kakebo"
)

var (
	storeConfiguration = k.StoreConfiguration{
		ConnectionTimeout: 1 * time.Second,
		DbName:            "kakebo_test.db",
	}
)

func TestStore(t *testing.T) {
	db, err := k.Connection(storeConfiguration)
	k.Check(err)
	defer db.Close()

	setup(db)

	t.Run("FetchExpensesForWeek(week: 1)", func(t *testing.T) {
		if !testFetchExpensesForWeek(db, 1) {
			t.Fail()
		}
	})
	t.Run("FetchExpensesForWeek(week: 2)", func(t *testing.T) {
		if !testFetchExpensesForWeek(db, 2) {
			t.Fail()
		}
	})
	t.Run("FetchExpensesForWeek(week: 3)", func(t *testing.T) {
		if !testFetchExpensesForWeek(db, 3) {
			t.Fail()
		}
	})

	t.Run("StoreExpense", func(t *testing.T) {
		expense := k.Expense{
			ID:          1,
			Description: "Toothpaste",
			Category:    "survival",
			Amount:      5,
			Week:        1,
		}
		if !testStoreExpense(db, expense) {
			t.Fail()
		}
	})

	t.Run("StoreIncome", func(t *testing.T) {
		incomes := []struct {
			monthYear string
			amount    int
		}{
			{
				monthYear: "December 2017",
				amount:    2000,
			},
			{
				monthYear: "January 2018",
				amount:    0,
			},
			{
				monthYear: "February 2018",
				amount:    10000,
			},
		}

		for _, tt := range incomes {
			if !testStoreIncome(db, tt.monthYear, tt.amount) {
				t.Fail()
			}
		}
	})

	t.Run("FetchIncome", func(t *testing.T) {
		incomes := []struct {
			monthYear string
			amount    int
		}{
			{
				monthYear: "December 2017",
				amount:    2000,
			},
			{
				monthYear: "January 2018",
				amount:    0,
			},
			{
				monthYear: "February 2018",
				amount:    10000,
			},
		}
		for _, income := range incomes {
			k.StoreIncome(db, income.monthYear, income.amount)
		}

		cases := []struct {
			monthYear string
			expected  int
		}{
			{
				monthYear: "December 2017",
				expected:  2000,
			},
			{
				monthYear: "October 2017",
				expected:  0,
			},
			{
				monthYear: "February 2018",
				expected:  10000,
			},
		}
		for _, tt := range cases {
			actual, _ := k.FetchIncome(db, tt.monthYear)
			if tt.expected != actual {
				t.Fail()
			}
		}
	})

	t.Run("StoreSavingsGoal", func(t *testing.T) {
		cases := []struct {
			monthYear string
			amount    string
		}{
			{
				monthYear: "January 2018",
				amount:    "1",
			},
			{
				monthYear: "December 2017",
				amount:    "2",
			},
			{
				monthYear: "September 2018",
				amount:    "3",
			},
		}

		for _, c := range cases {
			err := k.StoreSavingsGoal(db, c.monthYear, c.amount)
			if err != nil {
				t.Fail()
			}
		}
	})

	t.Run("FetchSavingsGoal", func(t *testing.T) {
		seeds := []struct {
			monthYear string
			amount    string
		}{
			{
				monthYear: "January 2018",
				amount:    "1000",
			},
			{
				monthYear: "February 2018",
				amount:    "2000",
			},
		}

		for _, s := range seeds {
			k.StoreSavingsGoal(db, s.monthYear, s.amount)
		}

		cases := []struct {
			monthYear string
			expected  int
		}{
			{
				monthYear: "January 2018",
				expected:  1000,
			},
			{
				monthYear: "February 2018",
				expected:  2000,
			},
			{
				monthYear: "October 2016",
				expected:  0,
			},
		}
		for _, c := range cases {
			actual, _ := k.FetchSavingsGoal(db, c.monthYear)
			if actual != c.expected {
				t.Fail()
			}
		}
	})

	teardown(db)
}

func testStoreIncome(db *bolt.DB, monthYear string, amount int) bool {
	err := k.StoreIncome(db, monthYear, amount)
	if err != nil {
		return false
	}
	return true
}

func testFetchExpensesForWeek(db *bolt.DB, week int) bool {
	expenses, err := k.FetchExpensesForWeek(db, week)
	k.Check(err)
	for _, expense := range expenses {
		if expense.Week != week {
			return false
		}
	}

	return true
}

func testStoreExpense(db *bolt.DB, expense k.Expense) bool {
	err := k.StoreExpense(db, expense)
	if err != nil {
		return false
	}

	return true
}

func seedExpenses(db *bolt.DB) {
	expenses := []k.Expense{
		{
			ID:          1,
			Description: "Toothpaste",
			Category:    "survival",
			Amount:      5,
			Week:        1,
		},
		{
			ID:          2,
			Description: "Vacation",
			Category:    "optional",
			Amount:      1000,
			Week:        1,
		},
		{
			ID:          3,
			Description: "Food",
			Category:    "survival",
			Amount:      17,
			Week:        2,
		},
		{
			ID:          4,
			Description: "Concert tickets",
			Category:    "extra",
			Amount:      120,
			Week:        2,
		},
	}

	for _, expense := range expenses {
		k.StoreExpense(db, expense)
	}
}

func setup(db *bolt.DB) {
	seedExpenses(db)
}

func teardown(db *bolt.DB) {
	err := k.ClearStore(db)
	k.Check(err)
}
