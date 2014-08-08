package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"appengine"
	"github.com/mjibson/appstats"
	"net/http"
)

func init() {
	http.Handle("/users", appstats.NewHandler(users))
}

func users(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "Accounts Listing"

	accounts, err := models.GetAllAccounts(c)

	if err != nil {
		page.Error = err.Error()
	} else {
		page.Accounts = accounts
	}

	pageTemplate := BuildTemplate(USERS)
	pageTemplate.Execute(w, page)
}
