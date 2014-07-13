package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"net/http"
)

func init() {
	http.HandleFunc("/myaccount", myaccount)
}

func myaccount(w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "My Account"
	pageTemplate := BuildTemplate(ACCOUNT)
	pageTemplate.Execute(w, page)
}
