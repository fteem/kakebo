package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

const (
	incomesBucket  = "incomes"
	expensesBucket = "expenses"
	savingsBucket  = "savings"
)

type Store struct {
	db   *bolt.DB
	path string
}

func NewStore(path string) *Store {
	return &Store{
		path: path,
	}
}

func (s *Store) Path() string { return s.path }

func (s *Store) Open() error {
	// Open database connection
	db, err := bolt.Open(s.path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	// Assign connection handler to Store
	s.db = db

	// Initialize needed buckets (if non-existent)
	if err := s.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(incomesBucket))
		tx.CreateBucketIfNotExists([]byte(expensesBucket))
		tx.CreateBucketIfNotExists([]byte(savingsBucket))
		return nil
	}); err != nil {
		s.Close()
		return err
	}

	return nil
}

func (s *Store) Close() error {
	if s.db != nil {
		s.db.Close()
	}

	return nil
}

func (s *Store) FetchExpensesForWeek(week int) ([]Expense, error) {
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

func (s *Store) FetchExpenses() ([]Expense, error) {
	var expenses []Expense

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(expensesBucket))

		c := b.Cursor()

		var expense Expense

		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &expense)
			if err != nil {
				return err
			}
			expenses = append(expenses, expense)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func (s *Store) StoreExpense(expense Expense) error {
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

func (s *Store) StoreIncome(monthYear string, amount int) error {
	incomeAsString := strconv.Itoa(amount)
	s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(incomesBucket))
		err := b.Put([]byte(monthYear), []byte(incomeAsString))
		return err
	})
	return nil
}

func (s *Store) FetchIncome(key string) (int, error) {
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

func (s *Store) FetchIncomes() ([]int, error) {
	var incomes []int
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(incomesBucket))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			amount, _ := strconv.Atoi(string(v))
			incomes = append(incomes, amount)
		}
		return nil
	})

	if err != nil {
		return []int{}, err
	}

	return incomes, nil
}

func (s *Store) StoreSavingsGoal(monthYear string, amount string) error {
	s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(savingsBucket))
		err := b.Put([]byte(monthYear), []byte(amount))
		return err
	})
	return nil
}

func (s *Store) FetchSavingsGoal(monthYear string) (int, error) {
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

func (s *Store) Clear() error {
	if err := s.db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(incomesBucket))
		tx.DeleteBucket([]byte(expensesBucket))
		tx.DeleteBucket([]byte(savingsBucket))
		return nil
	}); err != nil {
		return err
	}
	return nil
}
