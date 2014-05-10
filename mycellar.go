package pages

import (
	"models"
	"net/http"
)

func init() {
	http.HandleFunc("/mycellar", mycellar)
}

func mycellar(w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "My Cellar"

	pageTemplate := BuildTemplate(MYCELLAR)
	pageTemplate.Execute(w, page)
}
