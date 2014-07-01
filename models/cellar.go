package models

import (
	"fmt"
)

type Cellar struct {
	ID int
	NextBeerID int
	Name  string
	Beers []*Beer
}

type Beer struct {
	Name     string
	Notes    string
	Brewed   *Date
	Added    *Date
	Consumed []*Date
	Quantity int
}

func (cellar Cellar) addBeer(name, notes, brewed, added string, quantity int) {
	beer := newBeer(name, notes, brewed, added, quantity)
	cellar.Beers = append(cellar.Beers, beer)
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
		Name:   name,
		Notes:  notes,
		Brewed: brewedDate,
		Added:  addedDate,
		Consumed: []*Date{},
		Quantity: quantity,
	}
}

func (cellar *Cellar)AddBeer(beer *Beer) {
	cellar.Beers = append(cellar.Beers, beer)
}

func (beer *Beer)GetConsumedString() string {
	beer.Consumed = append(beer.Consumed, Now())
	beer.Consumed = append(beer.Consumed, Now())
	beer.Consumed = append(beer.Consumed, Now())
	str := ""
	for index, date := range beer.Consumed {
		if(index == 0) {
			str = date.ToString()
		} else {
			str = fmt.Sprintf("%s; %s", str, date.ToString())
		}
	}
	return str
}

func (beer *Beer)GetAgeString() string {
	startDate := beer.Brewed
	if(startDate == nil) {
		startDate = beer.Added
	}
	today := Now()
	difference := today.Date.Sub(startDate.Date)

	days := int(difference.Hours() / 24)
	months := int(days / (365/ 12))
	years := int(months / 12)

	days -= int(months * (365/ 12))
	months -= int(years * 12)

	ageStr := ""

	if(years > 0) {
		if(years > 1) {
			ageStr = fmt.Sprintf("%s %d years", ageStr, years)
		} else {
			ageStr = fmt.Sprintf("%s %d year", ageStr, years)
		}
	}
	if(months > 0) {
		if(months > 1) {
			ageStr = fmt.Sprintf("%s %d months", ageStr, months)
		} else {
			ageStr = fmt.Sprintf("%s %d month", ageStr, months)
		}
	}
	if(days > 0) {
		if(days > 1) {
			ageStr = fmt.Sprintf("%s %d days", ageStr, days)	
		} else {
			ageStr = fmt.Sprintf("%s %d day", ageStr, days)	
		}
		
	}
	return ageStr
}