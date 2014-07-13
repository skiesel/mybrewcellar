package models

import (
	"errors"
	"net/http"
	"strconv"
)

type Account struct {
	User         *User
	NextCellarID int
	Cellars      map[string]*Cellar
	CellarsByID  map[int]*Cellar
}

func (account Account) GetUsername() string {
	return account.User.UserID
}

func GuestAccount() *Account {
	return &Account{
		User: &User{
			UserID: "Guest",
			Email:  "",
		},
		NextCellarID: 0,
		Cellars:      map[string]*Cellar{},
	}
}

func NewAccount(userid, email string, r *http.Request) *Account {
	newAccount := &Account{
		User: &User{
			UserID: userid,
			Email:  email,
		},
		NextCellarID: 0,
		Cellars:      map[string]*Cellar{},
		CellarsByID:  map[int]*Cellar{},
	}

	newAccount.AddCellar("Long Term Storage")
	newAccount.AddCellar("Refrigerator")

	// c := appengine.NewContext(r)
	// k := datastore.NewKey(c, "Entity", email, 0, nil)
	// if _, err := datastore.Put(c, k, newAccount.toAccountDS()); err != nil {
	//  	c.Infof(err.Error())
	//  }

	return newAccount
}

func (account *Account) AddCellar(cellarName string) error {
	if account.Cellars[cellarName] != nil {
		return errors.New("Cellar Already Exists")
	}

	account.Cellars[cellarName] = &Cellar{
		ID:         account.NextCellarID,
		Name:       cellarName,
		NextBeerID: 0,
		Beers:      map[string]*Beer{},
		BeersByID:  map[int]*Beer{},
	}

	account.CellarsByID[account.NextCellarID] = account.Cellars[cellarName]
	account.NextCellarID++

	return nil
}

func (account *Account) GetCellarByID(idStr string) *Cellar {
	id, _ := strconv.Atoi(idStr)
	return account.CellarsByID[id]
}

func (account *Account) DeleteCellar(cellarName string) error {
	cellar := account.Cellars[cellarName]
	if cellar == nil {
		return errors.New("Cellar does not exist")
	}

	delete(account.CellarsByID, cellar.ID)
	delete(account.Cellars, cellarName)

	return nil
}
