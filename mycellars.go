package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"appengine"
	"github.com/mjibson/appstats"
	"net/http"
)

func init() {
	http.Handle("/mycellars", appstats.NewHandler(mycellars))
}

func mycellars(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "My Cellars"
	pageTemplate := BuildTemplate(MYCELLARS)
	pageTemplate.Execute(w, page)
}
