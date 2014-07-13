package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"net/http"
)

func init() {
	http.HandleFunc("/", index)
}

func index(w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "Home"

	pageTemplate := BuildTemplate(INDEX)
	pageTemplate.Execute(w, page)
}
