package main

import (
	"encoding/binary"
	"strconv"
	"time"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func CurrentWeekAsInt() int {
	_, week := time.Now().ISOWeek()

	return week
}

func CurrentWeekAsString() string {
	s := strconv.Itoa(CurrentWeekAsInt())

	return s
}

func CurrentMonth() string {
	now := time.Now()
	return now.Month().String()
}

func CurrentYear() string {
	now := time.Now()
	return strconv.Itoa(now.Year())
}

func MonthAndYear() string {
	now := time.Now()
	monthAndYear := now.Month().String() + ", " + strconv.Itoa(now.Year())

	return monthAndYear
}
