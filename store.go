package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

const (
	incomesBucket     = "incomes"
	expensesBucket    = "expenses"
	savingsBucket     = "savings"
	connectionTimeout = 1 * time.Second
)

type StoreConfiguration struct {
	ConnectionTimeout time.Duration
	DbName            string
}

type Store struct {
	db *bolt.DB
}

func NewStore(dbName string) Store {
	// Open database connection
	db, _ := bolt.Open(dbName, 0600, &bolt.Options{Timeout: connectionTimeout})

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

	return Store{db: db}
}

func (s Store) FetchExpensesForWeek(week int) ([]Expense, error) {
	var expenses []Expense
	if week == 0 {
		week = CurrentWeekAsInt()
	}

	err := s.db.View(func(tx *bolt.Tx) error {
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

func (s Store) StoreExpense(expense Expense) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(expensesBucket))

		id, _ := b.NextSequence()
		expense.ID = int(id)
		jsonBlob, err := json.Marshal(expense)
		if err != nil {
			return err
		}

		return b.Put(Itob(expense.ID), jsonBlob)
	})
}

func (s Store) StoreIncome(monthYear string, amount int) error {
	incomeAsString := strconv.Itoa(amount)
	s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(incomesBucket))
		err := b.Put([]byte(monthYear), []byte(incomeAsString))
		return err
	})
	return nil
}

func (s Store) FetchIncome(key string) (int, error) {
	var income int
	err := s.db.View(func(tx *bolt.Tx) error {
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

func (s Store) StoreSavingsGoal(monthYear string, amount string) error {
	s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(savingsBucket))
		err := b.Put([]byte(monthYear), []byte(amount))
		return err
	})
	return nil
}

func (s Store) FetchSavingsGoal(monthYear string) (int, error) {
	var savingsGoal int
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(savingsBucket))
		v := b.Get([]byte(monthYear))
		savingsGoal, _ = strconv.Atoi(string(v))
		return nil
	})

	if err != nil {
		return 0, err
	}
	return savingsGoal, nil
}

func (s Store) Clear() error {
	s.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(incomesBucket))
		if err != nil {
			return err
		}
		return nil
	})
	s.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(expensesBucket))
		if err != nil {
			return err
		}
		return nil
	})
	s.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(savingsBucket))
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (s Store) Close() {
	s.db.Close()
}
