package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"appengine"
	"appengine/datastore"
	"fmt"
	"github.com/mjibson/appstats"
	"net/http"
	"strconv"
)

func init() {
	http.Handle("/beer", appstats.NewHandler(beer))
	http.Handle("/universal-beer", appstats.NewHandler(universalbeer))
}

func beer(c appengine.Context, w http.ResponseWriter, r *http.Request) {
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

func universalbeer(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Redirect(w, r, "/mycellars", 303) //303 == See Other
	}

	q := datastore.NewQuery("Beer").Filter("UBID =", id)
	var beers []models.BeerDS
	beerKeys, err := q.GetAll(c, &beers)

	if len(beers) <= 0 {
		http.Redirect(w, r, "/mycellars", 303) //303 == See Other
	}

	cellar := &models.CellarDS{}
	cellarKey := beerKeys[0].Parent()
	datastore.Get(c, cellarKey, cellar)

	url := fmt.Sprintf("/beer?cellar=%d&id=%d", cellar.ID, beers[0].ID)
	http.Redirect(w, r, url, 303) //303 == See Other
}
