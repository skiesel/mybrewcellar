package models

import (
	"strconv"
	"strings"
	"time"
)

type Date struct {
	date time.Time
}

func Now() *Date {
	now := time.Now()
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil
	}
	time := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utc)
	return &Date{
		date: time,
	}
}

func ParseDate(date string) *Date {
	tokens := strings.Split(date, "-")
	if len(tokens) != 3 {
		return nil
	}
	year, err := strconv.Atoi(tokens[0])
	if err != nil {
		return nil
	}
	month, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil
	}
	day, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil
	}
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil
	}
	time := time.Date(year, time.Month(month), day, 0, 0, 0, 0, utc)
	return &Date{
		date: time,
	}
}
