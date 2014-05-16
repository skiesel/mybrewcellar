package models

import (
	"errors"
)

type Account struct {
	user    *User
	cellars map[string]*Cellar
}

func (account Account) GetUsername() string {
	return account.user.userid
}

func GuestAccount() *Account {
	return &Account{
		user: &User{
			userid: "Guest",
			email:  "",
		},
		cellars: map[string]*Cellar{},
	}
}

func NewAccount(userid, email string) *Account {
	return &Account{
		user: &User{
			userid: userid,
			email:  email,
		},
		cellars: map[string]*Cellar{
			"Long Term Storage": &Cellar{
				name:  "Long Term Storage",
				beers: []*Beer{},
			},
			"Refrigerator": &Cellar{
				name:  "Refrigerator",
				beers: []*Beer{},
			},
		},
	}
}

func GetAccount(email string) *Account {
	return nil
}

func (account Account) AddCellar(cellarName string) error {
	if account.cellars[cellarName] != nil {
		return errors.New("Cellar Already Exists")
	}

	account.cellars[cellarName] = &Cellar{
		name:  cellarName,
		beers: []*Beer{},
	}

	return nil
}

func (account Account) DeleteCellar(cellarName string) {
	delete(account.cellars, cellarName)
}

