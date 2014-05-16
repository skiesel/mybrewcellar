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

	account := models.GetAccount("test@test.com")
	if account == nil {
		account = models.NewAccount("username", "test@test.com")
	}

	account.DeleteCellar("NOT A CELLAR")
}
