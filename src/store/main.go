package store

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fteem/kakebo/src/util"
)

// Names
const (
	dbName         = "kakebo.db"
	incomesBucket  = "incomes"
	expensesBucket = "expenses"
	savingsBucket  = "savings"
)

// Settings
const (
	connectionTimeout = 1 * time.Second
)

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

func Connection() (*bolt.DB, error) {
	// Open database connection
	db, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: connectionTimeout})

	if err != nil {
		return nil, err
	}

	// Create config bucket if not present
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(incomesBucket))
		if err != nil {
			return err
		}
		return nil
	})

	// Create expenses bucket if not present
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(expensesBucket))
		if err != nil {
			return err
		}
		return nil
	})

	// Create savings bucket if not present
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(savingsBucket))
		if err != nil {
			return err
		}
		return nil
	})

	return db, nil
}

func FetchExpensesForWeek(db *bolt.DB, week int) ([]Expense, error) {
	var expenses []Expense
	if week == 0 {
		week = util.CurrentWeekAsInt()
	}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(expensesBucket))

		c := b.Cursor()

		var expense Expense

		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &expense)
			if err != nil {
				return err
			}
			if expense.Week == week {
				expenses = append(expenses, expense)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func StoreExpense(db *bolt.DB, expense Expense) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(expensesBucket))

		id, _ := b.NextSequence()
		expense.ID = int(id)
		jsonBlob, err := json.Marshal(expense)
		if err != nil {
			return err
		}

		return b.Put(util.Itob(expense.ID), jsonBlob)
	})
}

func StoreIncome(db *bolt.DB, key string, amount int) error {
	incomeAsString := strconv.Itoa(amount)
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(incomesBucket))
		err := b.Put([]byte(key), []byte(incomeAsString))
		return err
	})
	return nil
}

func FetchIncome(db *bolt.DB, key string) (int, error) {
	var income int
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(incomesBucket))
		v := b.Get([]byte(key))
		income, _ = strconv.Atoi(string(v))
		return nil
	})

	if err != nil {
		return 0, err
	}

	return income, nil
}

func StoreSavingsGoal(db *bolt.DB, key string, amount string) error {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(savingsBucket))
		err := b.Put([]byte(key), []byte(amount))
		return err
	})
	return nil
}

func FetchSavingsGoal(db *bolt.DB, key string) (int, error) {
	var savingsGoal int
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(savingsBucket))
		v := b.Get([]byte(key))
		savingsGoal, _ = strconv.Atoi(string(v))
		return nil
	})

	if err != nil {
		return 0, err
	}
	return savingsGoal, nil
}
