package models

import (
	"appengine"
	"appengine/user"
	"net/http"
)

type Page struct {
	Title string
	Username string
}

func NewPage(r *http.Request) Page {
	newPage := Page {
		Title : "Page Title",
		Username : "Guest",
	}

	c := appengine.NewContext(r)
	usr := user.Current(c)
	if usr != nil {
		newPage.Username = usr.String()
	}

	return newPage
}