package main_test

import (
	"testing"

	k "github.com/fteem/kakebo"
)

const (
	path = "kakebo_test.db"
)

func TestStore(t *testing.T) {
	store := k.NewStore(path)
	if err := store.Open(); err != nil {
		panic(err)
	}
	defer store.Close()

	setup(store)

	t.Run("FetchExpensesForWeek(month: January, year: 2018)", func(t *testing.T) {
		if !testFetchExpensesForMonth(store, "January", "2018") {
			t.Fail()
		}
	})
	t.Run("FetchExpensesForWeek(month: Ferbuary, year: 2017)", func(t *testing.T) {
		if !testFetchExpensesForMonth(store, "February", "2017") {
			t.Fail()
		}
	})
	t.Run("FetchExpensesForWeek(month: March, year: 2017)", func(t *testing.T) {
		if !testFetchExpensesForMonth(store, "March", "2017") {
			t.Fail()
		}
	})

	t.Run("StoreExpense", func(t *testing.T) {
		expense := k.Expense{
			ID:          1,
			Description: "Toothpaste",
			Category:    "survival",
			Amount:      5,
			Month:       "July",
			Year:        "2017",
		}
		if !testStoreExpense(store, expense) {
			t.Fail()
		}
	})

	t.Run("StoreIncome", func(t *testing.T) {
		incomes := []struct {
			month  string
			year   string
			amount int
		}{
			{
				month:  "December",
				year:   "2017",
				amount: 2000,
			},
			{
				month:  "January",
				year:   "2018",
				amount: 0,
			},
			{
				month:  "February",
				year:   "2018",
				amount: 10000,
			},
		}

		for _, tt := range incomes {
			if !testStoreIncome(store, tt.month, tt.year, tt.amount) {
				t.Fail()
			}
		}
	})

	t.Run("FetchIncome", func(t *testing.T) {
		incomes := []struct {
			month  string
			year   string
			amount int
		}{
			{
				month:  "December",
				year:   "2017",
				amount: 2000,
			},
			{
				month:  "January",
				year:   "2018",
				amount: 0,
			},
			{
				month:  "February",
				year:   "2018",
				amount: 10000,
			},
		}
		for _, income := range incomes {
			store.StoreIncome(income.month, income.year, income.amount)
		}

		for _, income := range incomes {
			actual, _ := store.FetchIncome(income.month, income.year)
			if income.amount != actual {
				t.Fail()
			}
		}
	})

	t.Run("StoreSavingsGoal", func(t *testing.T) {
		cases := []struct {
			month  string
			year   string
			amount string
		}{
			{
				month:  "January",
				year:   "2018",
				amount: "1",
			},
			{
				month:  "December",
				year:   "2017",
				amount: "2",
			},
			{
				month:  "September",
				year:   "2018",
				amount: "3",
			},
		}

		for _, c := range cases {
			err := store.StoreSavingsGoal(c.month, c.year, c.amount)
			if err != nil {
				t.Fail()
			}
		}
	})

	t.Run("FetchSavingsGoal", func(t *testing.T) {
		seeds := []struct {
			month  string
			year   string
			amount string
		}{
			{
				month:  "January",
				year:   "2018",
				amount: "1000",
			},
			{
				month:  "February",
				year:   "2018",
				amount: "2000",
			},
		}

		for _, s := range seeds {
			store.StoreSavingsGoal(s.month, s.year, s.amount)
		}

		cases := []struct {
			month    string
			year     string
			expected int
		}{
			{
				month:    "January",
				year:     "2018",
				expected: 1000,
			},
			{
				month:    "February",
				year:     "2018",
				expected: 2000,
			},
			{
				month:    "October",
				year:     "2016",
				expected: 0,
			},
		}
		for _, c := range cases {
			actual, _ := store.FetchSavingsGoal(c.month, c.year)
			if actual != c.expected {
				t.Fail()
			}
		}
	})

	teardown(store)
}

func testStoreIncome(store *k.Store, month string, year string, amount int) bool {
	err := store.StoreIncome(month, year, amount)
	if err != nil {
		return false
	}
	return true
}

func testFetchExpensesForMonth(store *k.Store, month string, year string) bool {
	expenses, err := store.FetchExpensesForMonth(month, year)
	k.Check(err)
	for _, expense := range expenses {
		if expense.Month != month && expense.Year != year {
			return false
		}
	}

	return true
}

func testStoreExpense(store *k.Store, expense k.Expense) bool {
	err := store.StoreExpense(expense)
	if err != nil {
		return false
	}

	return true
}

func seedExpenses(store *k.Store) {
	expenses := []k.Expense{
		{
			ID:          1,
			Description: "Toothpaste",
			Category:    "survival",
			Amount:      5,
			Month:       "June",
			Year:        "2017",
		},
		{
			ID:          2,
			Description: "Vacation",
			Category:    "optional",
			Amount:      1000,
			Month:       "July",
			Year:        "2017",
		},
		{
			ID:          3,
			Description: "Food",
			Category:    "survival",
			Amount:      17,
			Month:       "October",
			Year:        "2017",
		},
		{
			ID:          4,
			Description: "Concert tickets",
			Category:    "extra",
			Amount:      120,
			Month:       "March",
			Year:        "2017",
		},
	}

	for _, expense := range expenses {
		store.StoreExpense(expense)
	}
}

func setup(store *k.Store) {
	seedExpenses(store)
}

func teardown(store *k.Store) {
	err := store.Clear()
	k.Check(err)
}
