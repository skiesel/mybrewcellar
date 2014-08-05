package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"appengine"
	"github.com/mjibson/appstats"
	"net/http"
	"appengine/user"
	"encoding/csv"
	"io"
	"strconv"
)

func init() {
	http.Handle("/myaccount", appstats.NewHandler(myaccount))
	http.Handle("/export", appstats.NewHandler(exportaccount))
	http.Handle("/import", appstats.NewHandler(importaccount))
}

func myaccount(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	page := models.NewPage(r)
	page.Title = "My Account"
	pageTemplate := BuildTemplate(ACCOUNT)
	pageTemplate.Execute(w, page)
}

func exportaccount(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	usr := user.Current(c)
	if usr == nil {
		w.Write([]byte("Not logged in"))
		return
	}
	account := models.GetAccount(c, usr.Email)
	for _, cellar := range account.Cellars {
		w.Write([]byte(cellar.ToCSV()+"\n"))
		for _, beer := range cellar.Beers {
			w.Write([]byte(beer.ToCSV()+"\n"))
			for _, tasting := range beer.TastingsByID {
				w.Write([]byte(tasting.ToCSV()+"\n"))
			}
		}
	}
}

func importaccount(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	usr := user.Current(c)
	if usr == nil {
		w.Write([]byte("Not logged in"))
	}
	account := models.GetAccount(c, usr.Email)

	r.ParseForm()

	//file, header, err
	file, _, err := r.FormFile("importfile")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	reset := r.Form.Get("resetAccount") == "on"

	if reset {
		account.Cellars = map[string]*models.Cellar{}
		account.CellarsByID = map[int]*models.Cellar{}
	}
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	err = nil

	var currentCellar *models.Cellar
	var currentBeer *models.Beer

	for ; err == nil ; {
		line, err2 := csvReader.Read()
		err = err2
		if len(line) > 0 {
			switch {
				case "CELLAR" == line[0]: {
					w.Write([]byte("Got Cellar"))
					cellarName := line[1]
					currentCellar = account.Cellars[cellarName]
					if currentCellar == nil {
						currentCellar = &models.Cellar {
							ID : account.NextCellarID,
							NextBeerID : 0,
							Name : cellarName,
							Beers : map[string]*models.Beer{},
							BeersByID : map[int]*models.Beer{},
						}
						account.Cellars[cellarName] = currentCellar
						account.CellarsByID[account.NextCellarID] = currentCellar
						account.NextCellarID++
					}
					break
				}
				case "BEER" == line[0]: {
					w.Write([]byte("Got Beer"))
					ubid, err := strconv.Atoi(line[1])
					if err != nil {
						w.Write([]byte(err.Error()))
						return
					}
					if ubid < 0 {
						ubid, err = models.GetAndIncrementUniversalBeerID(c)
						if err != nil {
							w.Write([]byte(err.Error()))
							return
						}
					}
					beerName := line[2]
					quantity, err := strconv.Atoi(line[3])
					if err != nil {
						w.Write([]byte(err.Error()))
						return
					}
					notes := line[4]
					brewed := models.ParseDate(line[5])
					added := models.ParseDate(line[6])
					currentBeer = currentCellar.Beers[beerName]
					if currentBeer == nil {
						currentBeer = &models.Beer{
							UBID : ubid,
							ID : currentCellar.NextBeerID,
							Name : beerName,
							Notes : notes,
							Brewed : brewed,
							Added : added,
							Quantity : quantity,
							NextTastingID : 0,
							TastingsByID : map[int]*models.Tasting{},
						}
						currentCellar.Beers[beerName] = currentBeer
						currentCellar.BeersByID[currentCellar.NextBeerID] = currentBeer
						currentCellar.NextBeerID++
					}
					break
				}
				case "TASTING" == line[0]: {
					w.Write([]byte("Got Tasting"))
					rating, err := strconv.Atoi(line[1])
					if err != nil {
						w.Write([]byte(err.Error()))
						return
					}
					notes := line[2]
					tasted := models.ParseDate(line[3])
					var currentTasting *models.Tasting
					for _, tasting := range currentBeer.TastingsByID {
						if tasting.Rating == rating && tasting.Notes == notes && tasting.Date.ToString() == tasted.ToString() {
							currentTasting = tasting
							break
						}
					}
					if currentTasting == nil {
						currentTasting = &models.Tasting {
							ID : currentBeer.NextTastingID,
							Rating : rating,
							Notes : notes,
							Date : tasted,
						}
						currentBeer.TastingsByID[currentBeer.NextTastingID] = currentTasting
						currentBeer.NextTastingID++
					}
					break
				}
			}
		}
	}

	if err != io.EOF {
		w.Write([]byte(err.Error()))
		return
	}

	models.SaveAccount(c, account)
}
