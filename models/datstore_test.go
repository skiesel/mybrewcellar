package models

import (
	"appengine/aetest"
	"testing"
)

func getTestAccount(email string) *Account {
	account := &Account{
		User: &User{
			UserID: "TestAccount",
			Email:  email,
		},
		NextCellarID: 0,
		Cellars:      map[string]*Cellar{},
		CellarsByID:  map[int]*Cellar{},
	}

	cellar := &Cellar{
		ID:         0,
		Name:       "Test Cellar",
		NextBeerID: 1,
		Beers:      map[string]*Beer{},
		BeersByID:  map[int]*Beer{},
	}

	account.Cellars[cellar.Name] = cellar
	account.CellarsByID[cellar.ID] = cellar

	beer := &Beer{
		ID:            0,
		Name:          "Sierra Nevada Pale Ale",
		Notes:         "This is a delicious beer",
		Quantity:      10,
		Brewed:        ParseDate("2013-01-02"),
		Added:         Now(),
		NextTastingID: 1,
		Tastings:      []*Tasting{},
		TastingsByID:  map[int]*Tasting{},
	}

	cellar.Beers[beer.Name] = beer
	cellar.BeersByID[beer.ID] = beer

	tasting := &Tasting{
		ID:    0,
		Notes: "This beer tasted great today",
		Date:  Now(),
	}

	beer.Tastings = append(beer.Tastings, tasting)
	beer.TastingsByID[tasting.ID] = tasting

	return account
}

func ATestWriteToDataStore(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	defer c.Close()

	email := "test@test.com"

	account := getTestAccount(email)

	err = SaveAccount(c, account)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	account2 := GetAccount(c, email)

	if account2 == nil {
		t.Log("Did not find account")
		t.FailNow()
	}

	if !account.isEqual(account2, t) {
		t.Log("Accounts not equal")
		t.FailNow()
	}
}

func TestWriteToDataStoreWithDelete(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	defer c.Close()

	email := "test@test.com"

	account := getTestAccount(email)

	err = SaveAccount(c, account)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	beer := account.CellarsByID[0].BeersByID[0]
	delete(account.CellarsByID[0].Beers, beer.Name)
	delete(account.CellarsByID[0].BeersByID, 0)

	err = SaveAccount(c, account)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	account2 := GetAccount(c, email)

	if account2.CellarsByID[0] != account2.Cellars[account.CellarsByID[0].Name] {
		t.Log("Cellar pointers were different")
		t.FailNow()
	}

	if account2 == nil {
		t.Log("Did not find account")
		t.FailNow()
	}

	if !account.isEqual(account2, t) {
		t.Log("Accounts not equal")
		t.FailNow()
	}
}

func (a *Account) isEqual(b *Account, t *testing.T) bool {
	t.Log("...comparing accounts")
	if (a == nil && b != nil) || (a != nil && b == nil) {
		t.Log("One account was nil")
		return false
	} else if a == nil && b == nil {
		return true
	}
	match := a.User.UserID == b.User.UserID &&
		a.User.Email == b.User.Email &&
		a.NextCellarID == b.NextCellarID
	if !match {
		return false
	}

	if len(a.Cellars) != len(b.Cellars) {
		t.Log("Cellars were different sizes")
		return false
	}

	for id, cellar := range a.CellarsByID {
		if !cellar.isEqual(b.CellarsByID[id], t) {
			return false
		}
	}
	return true
}

func (a *Cellar) isEqual(b *Cellar, t *testing.T) bool {
	t.Log("...comparing cellars")
	if (a == nil && b != nil) || (a != nil && b == nil) {
		t.Log("One cellar was nil")
		return false
	} else if a == nil && b == nil {
		return true
	}

	match := a.ID == b.ID &&
		a.NextBeerID == b.NextBeerID &&
		a.Name == b.Name
	if !match {
		return false
	}

	if len(a.Beers) != len(b.Beers) {
		t.Log("Beers were different sizes")
		return false
	}

	for id, beer := range a.BeersByID {
		if !beer.isEqual(b.BeersByID[id], t) {
			return false
		}
	}

	return true
}

func (a *Beer) isEqual(b *Beer, t *testing.T) bool {
	t.Log("...comparing beers")
	if (a == nil && b != nil) || (a != nil && b == nil) {
		t.Log("One beer was nil")
		return false
	} else if a == nil && b == nil {
		return true
	}
	match := a.ID == b.ID &&
		a.Name == b.Name &&
		a.Notes == b.Notes &&
		a.Brewed.isEqual(b.Brewed, t) &&
		a.Added.isEqual(b.Added, t) &&
		a.Quantity == b.Quantity &&
		a.NextTastingID == b.NextTastingID
	if !match {
		return false
	}

	if len(a.Tastings) != len(b.Tastings) {
		t.Log("Tastings were different sizes")
		return false
	}

	for id, tasting := range a.TastingsByID {
		if !tasting.isEqual(b.TastingsByID[id], t) {
			return false
		}
	}

	return true
}

func (a *Tasting) isEqual(b *Tasting, t *testing.T) bool {
	t.Log("...comparing tastings")
	if (a == nil && b != nil) || (a != nil && b == nil) {
		t.Log("One tasting was nil")
		return false
	} else if a == nil && b == nil {
		return true
	}
	return a.ID == b.ID &&
		a.Notes == b.Notes &&
		a.Date.isEqual(b.Date, t)
}

func (a *Date) isEqual(b *Date, t *testing.T) bool {
	t.Log("...comparing dates")
	if (a == nil && b != nil) || (a != nil && b == nil) {
		t.Log("One date was nil")
		return false
	} else if a == nil && b == nil {
		return true
	}

	return a.Date == b.Date
}
