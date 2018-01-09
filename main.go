package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/fteem/kakebo/src/store"
	"github.com/fteem/kakebo/src/util"
)

var (
	app = kingpin.New("kakebo", "A household finance ledger")

	// Income
	income          = app.Command("income", "Income operations")
	incomeShow      = income.Command("show", "Show income")
	incomeSet       = income.Command("set", "Set income")
	incomeSetAmount = incomeSet.Arg("amount", "Income amount").Required().Int()

	// Savings
	savings = app.Command("savings", "Savings operations")

	// Target
	target          = app.Command("target", "Savings target operations")
	targetShow      = target.Command("show", "Show savings target")
	targetSet       = target.Command("set", "Set savings target")
	targetSetAmount = targetSet.Arg("amount", "Target amount").Required().String()

	// Expenses
	expenses = app.Command("expenses", "Expenses operations")

	expensesAdd            = expenses.Command("add", "Add expense")
	expensesAddDescription = expensesAdd.Flag("description", "Expense description").Short('d').Required().String()
	expensesAddAmount      = expensesAdd.Flag("amount", "Expense amount").Short('a').Required().Int()
	expensesAddCategory    = expensesAdd.Flag("category", "Expense category").Short('c').Required().Enum("survival", "optional", "culture", "extra")
	expensesAddWeek        = expensesAdd.Flag("week", "Week number").Short('w').Default(util.CurrentWeekAsString()).Int()

	expensesList     = expenses.Command("list", "List expenses")
	expensesListWeek = expensesList.Flag("week", "Week number").Short('w').Default(util.CurrentWeekAsString()).Int()
)

func main() {
	db, err := store.Connection()
	util.Check(err)
	defer db.Close()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case incomeShow.FullCommand():
		income, err := store.FetchIncome(db, util.MonthAndYear())
		util.Check(err)
		fmt.Println(income)
	case incomeSet.FullCommand():
		err := store.StoreIncome(db, util.MonthAndYear(), *incomeSetAmount)
		util.Check(err)
	case targetSet.FullCommand():
		err := store.StoreSavingsGoal(db, util.MonthAndYear(), *targetSetAmount)
		util.Check(err)
	case targetShow.FullCommand():
		goal, err := store.FetchSavingsGoal(db, util.MonthAndYear())
		util.Check(err)
		fmt.Println("This month's goal:", goal)
	case expensesAdd.FullCommand():
		expense := store.Expense{
			Description: *expensesAddDescription,
			Amount:      *expensesAddAmount,
			Category:    *expensesAddCategory,
			Week:        *expensesAddWeek,
		}
		err := store.StoreExpense(db, expense)
		util.Check(err)
	case expensesList.FullCommand():
		expenses, err := store.FetchExpensesForWeek(db, *expensesListWeek)
		util.Check(err)

		for _, expense := range expenses {
			fmt.Println(expense)
		}
	}
}
