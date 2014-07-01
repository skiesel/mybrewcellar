package pages

import (
	"models"
	"net/http"
)

func init() {
	http.HandleFunc("/cellar", cellar)
}

func cellar(w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "Cellar"

	id := r.URL.Query().Get("id")
	if(id == "") {
		http.Redirect(w, r, "/mycellars", 303) //303 == See Other
	}
	
	page.Cellar = page.Account.GetCellarById(id)

	pageTemplate := BuildTemplate(CELLAR)
	pageTemplate.Execute(w, page)
}
