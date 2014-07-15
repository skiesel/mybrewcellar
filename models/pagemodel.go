package models

import (
	"appengine"
	"appengine/user"
	"net/http"
)

type Page struct {
	Title   string
	Logout  string
	Account *Account
	Cellar  *Cellar
	Beer    *Beer
	Error   string
}

func NewPage(r *http.Request) Page {
	newPage := Page{
		Title:   "Page Title",
		Logout:  "",
		Account: GuestAccount(),
		Cellar:  nil,
		Beer:    nil,
		Error:   "",
	}

	c := appengine.NewContext(r)
	usr := user.Current(c)

	if usr != nil {
		logout, _ := user.LogoutURL(c, "/")
		newPage.Logout = logout
		newPage.Account = GetAccount(c, usr.Email)

		if newPage.Account == nil {
			newPage.Account = NewAccount("NewUser", usr.Email, r)
			beer := newBeer("Sierra Nevada Pale Ale", "", "2013-07-01", "", 1)
			beer.ID = 0
			newPage.Account.CellarsByID[0].Beers[beer.Name] = beer
			newPage.Account.CellarsByID[0].BeersByID[beer.ID] = beer

			beer = newBeer("Sierra Nevada Bigfoot Barleywine", "", "2014-06-01", "", 10)
			beer.ID = 1
			beer.TastingsByID[0] = &Tasting{
				ID : 0,
				Rating : 5,
				Notes: "Sample Note",
				Date:  Now(),
			}

			beer.NextTastingID = 1

			newPage.Account.CellarsByID[0].Beers[beer.Name] = beer
			newPage.Account.CellarsByID[0].BeersByID[beer.ID] = beer

		newPage.Account.CellarsByID[0].NextBeerID = 2

			SaveAccount(c, newPage.Account)
		}
	}
	return newPage
}
