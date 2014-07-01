package models

import (
	"errors"
	"strconv"
)

type Account struct {
	User    *User
	NextCellarID int
	Cellars map[string]*Cellar
	CellarsByID map[int]*Cellar
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
		Cellars: map[string]*Cellar{},
	}
}

func NewAccount(userid, email string) *Account {
	newAccount := &Account{
		User: &User{
			UserID: userid,
			Email:  email,
		},
		NextCellarID: 0,
		Cellars: map[string]*Cellar{},
		CellarsByID: map[int]*Cellar{},
	}

	newAccount.AddCellar("Long Term Storage")
	newAccount.AddCellar("Refrigerator")

	return newAccount
}

func GetAccount(email string) *Account {
	return nil
}

func (account *Account) AddCellar(cellarName string) error {
	if account.Cellars[cellarName] != nil {
		return errors.New("Cellar Already Exists")
	}

	account.Cellars[cellarName] = &Cellar{
		ID: account.NextCellarID,
		Name:  cellarName,
		NextBeerID: 0,
		Beers: []*Beer{},
	}

	account.CellarsByID[account.NextCellarID] = account.Cellars[cellarName]
	account.NextCellarID++

	return nil
}

func (account Account) GetCellarById(idStr string) *Cellar {
	id, _ := strconv.Atoi(idStr)
	return account.CellarsByID[id]
}

func (account Account) DeleteCellar(cellarName string) {
	cellar, exists := account.Cellars[cellarName]
	if(exists) {
		delete(account.CellarsByID, cellar.ID)
		delete(account.Cellars, cellarName)
	}
}

