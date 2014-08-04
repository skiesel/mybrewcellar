package models

import (
	"strconv"
)

type Cellar struct {
	ID         int
	NextBeerID int
	Name       string
	Beers      map[string]*Beer
	BeersByID  map[int]*Beer
}

type CellarDS struct {
	ID         int
	NextBeerID int
	Name       string
}

func (cellar Cellar) toCellarDS() *CellarDS {
	return &CellarDS{
		ID:         cellar.ID,
		NextBeerID: cellar.NextBeerID,
		Name:       cellar.Name,
	}
}

func (cellarDS *CellarDS) toCellar() *Cellar {
	return &Cellar{
		ID:         cellarDS.ID,
		NextBeerID: cellarDS.NextBeerID,
		Name:       cellarDS.Name,
		Beers:      map[string]*Beer{},
		BeersByID:  map[int]*Beer{},
	}
}

func (cellar *Cellar) GetBeerByID(idStr string) *Beer {
	id, _ := strconv.Atoi(idStr)
	return cellar.BeersByID[id]
}
