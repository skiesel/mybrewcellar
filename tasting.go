package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"appengine"
	"github.com/mjibson/appstats"
	"net/http"
)

func init() {
	http.Handle("/tasting", appstats.NewHandler(tasting))
}

func tasting(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "Tasting"

	cellar := r.URL.Query().Get("cellar")
	beer := r.URL.Query().Get("beer")
	id := r.URL.Query().Get("id")
	if id == "" || beer == "" || cellar == "" {
		http.Redirect(w, r, "/mycellars", 303) //303 == See Other
	}

	page.Cellar = page.Account.GetCellarByID(cellar)
	page.Beer = page.Cellar.GetBeerByID(id)
	page.Tasting = page.Beer.GetTastingByID(id)

	pageTemplate := BuildTemplate(TASTING)
	pageTemplate.Execute(w, page)
}
