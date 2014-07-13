package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"net/http"
)

func init() {
	http.HandleFunc("/beer", beer)
}

func beer(w http.ResponseWriter, r *http.Request) {
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
