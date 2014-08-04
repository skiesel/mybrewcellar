package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"appengine"
	"github.com/mjibson/appstats"
	"net/http"
)

func init() {
	http.Handle("/cellar", appstats.NewHandler(cellar))
}

func cellar(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "Cellar"

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/mycellars", 303) //303 == See Other
	}

	page.Cellar = page.Account.GetCellarByID(id)

	pageTemplate := BuildTemplate(CELLAR)
	pageTemplate.Execute(w, page)
}
