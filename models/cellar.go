package models

type Cellar struct {
	id int
	nextBeerId int
	name  string
	beers []*Beer
}

type Beer struct {
	name     string
	notes    string
	brewed   *Date
	added    *Date
	consumed *Date
}

func (cellar Cellar) addBeer(name, notes, brewed, added string) {
	beer := newBeer(name, notes, brewed, added)
	cellar.beers = append(cellar.beers, beer)
}

func newBeer(name, notes, brewed, added string) *Beer {
	var brewedDate *Date
	if brewed != "" {
		brewedDate = ParseDate(brewed)
	}

	addedDate := Now()
	if added != "" {
		addedDate = ParseDate(added)
	}

	return &Beer{
		name:   name,
		notes:  notes,
		brewed: brewedDate,
		added:  addedDate,
	}
}

func (cellar *Cellar)GetName() string {
	return cellar.name
}

func (cellar *Cellar)GetId() int {
	return cellar.id
}

func (cellar *Cellar)GetBeers() []*Beer {
	return cellar.beers
}

func (beer *Beer)GetName() string {
	return beer.name
}