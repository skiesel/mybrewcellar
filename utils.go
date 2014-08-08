package pages

import (
	"html/template"
)

const (
	INDEX     = "index.html"
	MYCELLARS = "mycellars.html"
	CELLAR    = "cellar.html"
	ACCOUNT   = "account.html"
	BEER      = "beer.html"
	TASTING   = "tasting.html"
	USERS   = "users.html"
)

func BuildTemplate(templateName string) *template.Template {
	pageTemplate, err := template.New(templateName).ParseFiles("templates/top.html", "templates/bottom.html", "templates/"+templateName)
	if err != nil {
		pageTemplate, err = template.New("Error").Parse("Oops! Didn't see that one coming.")
		if err != nil {
			panic(err)
		}
	}
	return pageTemplate
}
