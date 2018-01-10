package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	path = "kakebo.db"
)

var (
	app = kingpin.New("kakebo", "A household finance ledger")

	// Income
	income = app.Command("income", "Income operations")

	incomeShow      = income.Command("show", "Show income")
	incomeShowMonth = incomeShow.Flag("month", "Month").Short('m').Default(CurrentMonth()).String()
	incomeShowYear  = incomeShow.Flag("year", "Year").Short('y').Default(CurrentYear()).String()

	incomeSet       = income.Command("set", "Set income")
	incomeSetAmount = incomeSet.Arg("amount", "Income amount").Required().Int()
	incomeSetMonth  = incomeSet.Flag("month", "Month").Short('m').Default(CurrentMonth()).String()
	incomeSetYear   = incomeSet.Flag("year", "Year").Short('y').Default(CurrentYear()).String()

	// Savings
	savings      = app.Command("savings", "Savings operations")
	savingsMonth = savings.Flag("month", "Month").Short('m').Default(CurrentMonth()).String()
	savingsYear  = savings.Flag("year", "Year").Short('y').Default(CurrentYear()).String()

	// Target
	target          = app.Command("target", "Savings target operations")
	targetShow      = target.Command("show", "Show savings target")
	targetShowMonth = targetShow.Flag("month", "Month").Short('m').Default(CurrentMonth()).String()
	targetShowYear  = targetShow.Flag("year", "Year").Short('y').Default(CurrentYear()).String()

	targetSet       = target.Command("set", "Set savings target")
	targetSetAmount = targetSet.Arg("amount", "Target amount").Required().String()
	targetSetMonth  = targetSet.Flag("month", "Month").Short('m').Default(CurrentMonth()).String()
	targetSetYear   = targetSet.Flag("year", "Year").Short('y').Default(CurrentYear()).String()

	// Expenses
	expenses = app.Command("expenses", "Expenses operations")

	expensesAdd            = expenses.Command("add", "Add expense")
	expensesAddDescription = expensesAdd.Flag("description", "Expense description").Short('d').Required().String()
	expensesAddAmount      = expensesAdd.Flag("amount", "Expense amount").Short('a').Required().Int()
	expensesAddCategory    = expensesAdd.Flag("category", "Expense category").Short('c').Required().Enum("survival", "optional", "culture", "extra")
	expensesAddMonth       = expensesAdd.Flag("month", "Month").Short('m').Default(CurrentMonth()).String()
	expensesAddYear        = expensesAdd.Flag("year", "Year").Short('y').Default(CurrentYear()).String()

	expensesList      = expenses.Command("list", "List expenses")
	expensesListMonth = expensesList.Flag("month", "Month").Short('m').Default(CurrentMonth()).String()
	expensesListYear  = expensesList.Flag("yaer", "Year").Short('y').Default(CurrentYear()).String()
)

func main() {
	store := NewStore(path)
	if err := store.Open(); err != nil {
		fmt.Errorf("Open store: %s", err)
	}

	defer store.Close()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case incomeShow.FullCommand():
		income, err := store.FetchIncome(*incomeShowMonth, *incomeShowYear)
		Check(err)
		fmt.Println(income)
	case incomeSet.FullCommand():
		err := store.StoreIncome(*incomeSetMonth, *incomeSetYear, *incomeSetAmount)
		Check(err)
	case targetSet.FullCommand():
		err := store.StoreSavingsGoal(*targetSetMonth, *targetSetYear, *targetSetAmount)
		Check(err)
	case targetShow.FullCommand():
		goal, err := store.FetchSavingsGoal(*targetShowMonth, *targetShowYear)
		Check(err)
		fmt.Println("This month's goal:", goal)
	case expensesAdd.FullCommand():
		expense := Expense{
			Description: *expensesAddDescription,
			Amount:      *expensesAddAmount,
			Category:    *expensesAddCategory,
			Month:       *expensesAddMonth,
			Year:        *expensesAddYear,
		}
		err := store.StoreExpense(expense)
		Check(err)
	case expensesList.FullCommand():
		expenses, err := store.FetchExpensesForMonth(*expensesListMonth, *expensesListYear)
		Check(err)

		for _, expense := range expenses {
			fmt.Println(expense)
		}
	case savings.FullCommand():
		expenses, err := store.FetchExpensesForMonth(*savingsMonth, *savingsYear)
		Check(err)

		expensesSum := 0
		for i := 0; i < len(expenses); i++ {
			expensesSum += expenses[i].Amount
		}

		income, err := store.FetchIncome(*savingsMonth, *savingsYear)
		Check(err)

		goal, err := store.FetchSavingsGoal(*savingsMonth, *savingsYear)
		Check(err)

		fmt.Println("Savings goal:", goal, "Total savings:", income-expensesSum)
	}
}
