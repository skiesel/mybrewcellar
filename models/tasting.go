package models

type Tasting struct {
	ID     int
	Rating int
	Notes  string
	Date   *Date
}

type TastingDS struct {
	ID     int
	Rating int
	Notes  string
	Date   string
}

func (tasting *Tasting) toTastingDS() *TastingDS {
	return &TastingDS{
		ID:     tasting.ID,
		Rating: tasting.Rating,
		Notes:  tasting.Notes,
		Date:   tasting.Date.ToString(),
	}
}

func (tastingDS *TastingDS) toTasting() *Tasting {
	return &Tasting{
		ID:     tastingDS.ID,
		Rating: tastingDS.Rating,
		Notes:  tastingDS.Notes,
		Date:   ParseDate(tastingDS.Date),
	}
}
