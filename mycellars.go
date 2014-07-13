package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"net/http"
)

func init() {
	http.HandleFunc("/mycellars", mycellars)
}

func mycellars(w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "My Cellars"
	pageTemplate := BuildTemplate(MYCELLARS)
	pageTemplate.Execute(w, page)
}
