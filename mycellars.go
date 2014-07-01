package pages

import (
	"models"
	"net/http"
)

func init() {
	http.HandleFunc("/mycellars", mycellars)
}

func mycellars(w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "My Cellars"

	newCellar := r.PostFormValue("newcellar")

	if(newCellar != "") {
		err := page.Account.AddCellar(newCellar)
		if(err != nil) {
			page.Error = err.Error();
		}
	}

	pageTemplate := BuildTemplate(MYCELLARS)
	pageTemplate.Execute(w, page)
}
