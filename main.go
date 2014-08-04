package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"appengine"
	"github.com/mjibson/appstats"
	"net/http"
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
