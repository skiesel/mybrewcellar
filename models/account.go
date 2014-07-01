package models

import (
	"errors"
	"strconv"
)

type Account struct {
	user    *User
	nextCellarId int
	cellars map[string]*Cellar
	cellarsById map[int]*Cellar
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
		nextCellarId: 0,
		cellars: map[string]*Cellar{},
	}
}

func NewAccount(userid, email string) *Account {
	lts := &Cellar{
				id: 0,
				nextBeerId: 0,
				name:  "Long Term Storage",
				beers: []*Beer{},
			}
	fridge := &Cellar{
				id: 1,
				nextBeerId: 0,
				name:  "Refrigerator",
				beers: []*Beer{},
			}

	return &Account{
		user: &User{
			userid: userid,
			email:  email,
		},
		nextCellarId: 2,
		cellars: map[string]*Cellar{
			"Long Term Storage": lts,
			"Refrigerator": fridge,
		},
		cellarsById: map[int]*Cellar{
			0:lts,
			1:fridge,
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

	account.cellarsById[account.nextCellarId] = account.cellars[cellarName]
	account.nextCellarId++

	return nil
}

func (account Account) GetCellars() map[string]*Cellar {
	return account.cellars
}

func (account Account) GetCellarById(idStr string) *Cellar {
	id, _ := strconv.Atoi(idStr)
	return account.cellarsById[id]
}

func (account Account) DeleteCellar(cellarName string) {
	cellar, exists := account.cellars[cellarName]
	if(exists) {
		delete(account.cellarsById, cellar.id)
		delete(account.cellars, cellarName)
	}
}

