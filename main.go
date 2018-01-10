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
	expensesAddWeek        = expensesAdd.Flag("week", "Week number").Short('w').Default(CurrentWeekAsString()).Int()

	expensesList     = expenses.Command("list", "List expenses")
	expensesListWeek = expensesList.Flag("week", "Week number").Short('w').Default(CurrentWeekAsString()).Int()
)

func main() {
	store := NewStore(path)
	if err := store.Open(); err != nil {
		fmt.Errorf("Open store: %s", err)
	}

	defer store.Close()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case incomeShow.FullCommand():
		income, err := store.FetchIncome(MonthAndYear())
		Check(err)
		fmt.Println(income)
	case incomeSet.FullCommand():
		err := store.StoreIncome(MonthAndYear(), *incomeSetAmount)
		Check(err)
	case targetSet.FullCommand():
		err := store.StoreSavingsGoal(MonthAndYear(), *targetSetAmount)
		Check(err)
	case targetShow.FullCommand():
		goal, err := store.FetchSavingsGoal(MonthAndYear())
		Check(err)
		fmt.Println("This month's goal:", goal)
	case expensesAdd.FullCommand():
		expense := Expense{
			Description: *expensesAddDescription,
			Amount:      *expensesAddAmount,
			Category:    *expensesAddCategory,
			Week:        *expensesAddWeek,
		}
		err := store.StoreExpense(expense)
		Check(err)
	case expensesList.FullCommand():
		expenses, err := store.FetchExpensesForWeek(*expensesListWeek)
		Check(err)

		for _, expense := range expenses {
			fmt.Println(expense)
		}
	case savings.FullCommand():
		expenses, err := store.FetchExpenses()
		Check(err)

		incomes, err := store.FetchIncomes()
		Check(err)

		expensesSum := 0
		for i := 0; i < len(expenses); i++ {
			expensesSum += expenses[i].Amount
		}

		incomesSum := 0
		for i := 0; i < len(incomes); i++ {
			incomesSum += incomes[i]
		}

		fmt.Println("Total savings: ", incomesSum-expensesSum)
	}
}
