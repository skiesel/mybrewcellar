package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"net/http"
	"github.com/mjibson/appstats"
	"appengine"
)

func init() {
	http.Handle("/myaccount", appstats.NewHandler(myaccount))
}

func myaccount(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "My Account"
	pageTemplate := BuildTemplate(ACCOUNT)
	pageTemplate.Execute(w, page)
}
