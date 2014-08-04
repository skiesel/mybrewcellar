package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"appengine"
	"github.com/mjibson/appstats"
	"net/http"
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
