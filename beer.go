package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"appengine"
	"github.com/mjibson/appstats"
	"net/http"
)

func init() {
	http.Handle("/beer", appstats.NewHandler(beer))
}

func beer(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "Beer"

	cellar := r.URL.Query().Get("cellar")
	id := r.URL.Query().Get("id")
	if id == "" || cellar == "" {
		http.Redirect(w, r, "/mycellars", 303) //303 == See Other
	}

	page.Cellar = page.Account.GetCellarByID(cellar)
	page.Beer = page.Cellar.GetBeerByID(id)

	pageTemplate := BuildTemplate(BEER)
	pageTemplate.Execute(w, page)
}
