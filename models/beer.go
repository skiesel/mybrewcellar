package models

import (
	"strconv"
	"bytes"
	"encoding/csv"
)

type Beer struct {
	UBID          int
	ID            int
	Name          string
	Notes         string
	Brewed        *Date
	Added         *Date
	Quantity      int
	NextTastingID int
	TastingsByID  map[int]*Tasting
}

type BeerDS struct {
	UBID          int
	ID            int
	Name          string
	Notes         string
	Brewed        string
	Added         string
	Quantity      int
	NextTastingID int
}

func (beer *Beer) toBeerDS() *BeerDS {
	return &BeerDS{
		UBID:          beer.UBID,
		ID:            beer.ID,
		Name:          beer.Name,
		Notes:         beer.Notes,
		Brewed:        beer.Brewed.ToString(),
		Added:         beer.Added.ToString(),
		Quantity:      beer.Quantity,
		NextTastingID: beer.NextTastingID,
	}
}

func (beerDS *BeerDS) toBeer() *Beer {
	return &Beer{
		UBID:          beerDS.UBID,
		ID:            beerDS.ID,
		Name:          beerDS.Name,
		Notes:         beerDS.Notes,
		Brewed:        ParseDate(beerDS.Brewed),
		Added:         ParseDate(beerDS.Added),
		Quantity:      beerDS.Quantity,
		NextTastingID: beerDS.NextTastingID,
		TastingsByID:  map[int]*Tasting{},
	}
}

func (beer *Beer) GetBirthday() *Date {
	birthday := beer.Brewed
	if birthday == nil {
		birthday = beer.Added
	}
	return birthday
}

func (beer *Beer) GetAverageRating() float64 {
	average := 0.0
	tastingCount := len(beer.TastingsByID)
	if tastingCount > 0 {
		for _, tasting := range beer.TastingsByID {
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

func (beer *Beer) GetAgeString() string {
	birthday := beer.GetBirthday()
	today := Now()
	return getDurationString(birthday, today)
}

func (beer *Beer) GetTastingByID(idStr string) *Tasting {
	id, _ := strconv.Atoi(idStr)
	return beer.TastingsByID[id]
}

func (beer *Beer) ToCSV() string {
	buf := new(bytes.Buffer)
	csvWriter := csv.NewWriter(buf)
	csvWriter.Write([]string{"BEER", strconv.Itoa(beer.UBID),
		beer.Name, strconv.Itoa(beer.Quantity), beer.Notes, beer.Brewed.ToString(), beer.Added.ToString()})
	csvWriter.Flush()
	return buf.String()
}