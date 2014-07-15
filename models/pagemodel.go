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
			SaveAccount(c, newPage.Account)
		}
	}
	return newPage
}
