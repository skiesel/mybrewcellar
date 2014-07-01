package models

import (
	"appengine"
	"appengine/user"
	"net/http"
)

type Page struct {
	Title    string
	Logout	 string
	Account  *Account
	Cellar	 *Cellar
	Error 	 string
}

func NewPage(r *http.Request) Page {
	newPage := Page{
		Title : "Page Title",
		Logout : "",
		Account : GuestAccount(),
	}

	c := appengine.NewContext(r)
	usr := user.Current(c)

	if usr != nil {
		logout, _ := user.LogoutURL(c, "/")
		newPage.Logout = logout
		newPage.Account = GetAccount(usr.Email)
		if newPage.Account == nil {
			newPage.Account = NewAccount("NewUser", usr.Email)
			newPage.Account.CellarsByID[0].AddBeer(newBeer("Sierra Nevada Pale Ale", "", "2013-07-01", "", 1))
			newPage.Account.CellarsByID[0].AddBeer(newBeer("Sierra Nevada Bigfoot Barleywine", "", "2014-06-01", "", 10))
		}
	}
	return newPage
}
