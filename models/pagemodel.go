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
		}
	}

	return newPage
}
