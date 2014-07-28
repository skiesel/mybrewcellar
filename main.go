package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"net/http"
	"github.com/mjibson/appstats"
	"appengine"
)

func init() {
	http.Handle("/", appstats.NewHandler(index))
}

func index(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "Home"

	pageTemplate := BuildTemplate(INDEX)
	pageTemplate.Execute(w, page)
}
