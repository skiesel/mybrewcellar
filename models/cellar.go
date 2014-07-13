package models

import (
	"fmt"
	"strings"
	"strconv"
)

type Cellar struct {
	ID         int
	NextBeerID int
	Name       string
	Beers      map[string]*Beer
	BeersByID  map[int]*Beer
}

type Beer struct {
	ID            int
	Name          string
	Notes         string
	Brewed        *Date
	Added         *Date
	Quantity      int
	NextTastingID int
	Tastings      []*Tasting
	TastingsByID  map[int]*Tasting
}

type Tasting struct {
	ID    int
	Rating int
	Notes string
	Date  *Date
}

func (beer *Beer) GetBirthday() *Date {
	birthday := beer.Brewed
	if birthday == nil {
		birthday = beer.Added
	}
	return birthday
}

func newBeer(name, notes, brewed, added string, quantity int) *Beer {
	var brewedDate *Date
	if brewed != "" {
		brewedDate = ParseDate(brewed)
	}

	addedDate := Now()
	if added != "" {
		addedDate = ParseDate(added)
	}

	return &Beer{
		ID: 0,
		Name:     name,
		Notes:    notes,
		Brewed:   brewedDate,
		Added:    addedDate,
		Tastings: []*Tasting{},
		Quantity: quantity,
	}
}

func (cellar *Cellar) GetBeerByID(idStr string) *Beer {
	id, _ := strconv.Atoi(idStr)
	return cellar.BeersByID[id]
}

func (cellar *Cellar) AddBeer(beer *Beer) {
	cellar.Beers[beer.Name] = beer
	cellar.BeersByID[beer.ID] = beer
}

func (cellar Cellar) addBeer(name, notes, brewed, added string, quantity int) {
	cellar.AddBeer(newBeer(name, notes, brewed, added, quantity))
}

func (beer *Beer) GetAverageRating() float64 {
	average := 0.0
	tastingCount := len(beer.Tastings)
	if tastingCount > 0 {
		for _, tasting := range beer.Tastings {
			average += float64(tasting.Rating)
		}
		average /= float64(tastingCount)
	}
	return average
}

func (beer *Beer) GetTastingAge(tasting *Tasting) string {
	birthday := beer.GetBirthday()
	return getDurationString(birthday, tasting.Date)
}

func getDurationString(from, to *Date) string {
	difference := to.Date.Sub(from.Date)

	days := int(difference.Hours() / 24)
	months := int(days / (365 / 12))
	years := int(months / 12)

	days -= int(months * (365 / 12))
	months -= int(years * 12)

	ageStr := ""

	if years > 0 {
		if years > 1 {
			ageStr = fmt.Sprintf("%s %d years", ageStr, years)
		} else {
			ageStr = fmt.Sprintf("%s %d year", ageStr, years)
		}
	}
	if months > 0 {
		if months > 1 {
			ageStr = fmt.Sprintf("%s %d months", ageStr, months)
		} else {
			ageStr = fmt.Sprintf("%s %d month", ageStr, months)
		}
	}
	if days > 0 {
		if days > 1 {
			ageStr = fmt.Sprintf("%s %d days", ageStr, days)
		} else {
			ageStr = fmt.Sprintf("%s %d day", ageStr, days)
		}

	}
	return strings.TrimSpace(ageStr)
}

func (beer *Beer) GetAgeString() string {
	birthday := beer.GetBirthday()
	today := Now()
	return getDurationString(birthday, today)
}
