package models

import (
	"appengine"
	"appengine/user"
	"appengine/datastore"
	"net/http"
)

type Page struct {
	Title   string
	Logout  string
	Account *Account
	Cellar  *Cellar
	Beer    *Beer
	Tasting *Tasting
	Error   string
	Accounts []string
	Editable bool
}

func NewPage(r *http.Request) Page {
	newPage := Page{
		Title:   "Page Title",
		Logout:  "",
		Account: GuestAccount(),
		Cellar:  nil,
		Beer:    nil,
		Tasting: nil,
		Error:   "",
		Accounts: []string{},
		Editable: true,
	}

	c := appengine.NewContext(r)
	usr := user.Current(c)

	if usr != nil {
		logout, _ := user.LogoutURL(c, "/")
		newPage.Logout = logout
		
		username := r.URL.Query().Get("username")
		
		if username == "" {
			newPage.Account = GetAccount(c, usr.Email)
		} else {
			query := datastore.NewQuery("Account").Filter("UserID =", username)
			var accounts []AccountDS
      _, err := query.GetAll(c, &accounts)
      
      if err != nil || len(accounts) <= 0 {
      	newPage.Account = GetAccount(c, usr.Email)
    	} else {
    		newPage.Account = GetAccount(c, accounts[0].UserEmail)
    		newPage.Editable = false
    	}
			
		}

		if newPage.Account == nil {
			newPage.Account = NewAccount("NewUser", usr.Email, r)
			SaveAccount(c, newPage.Account)
		}

	}
	return newPage
}
