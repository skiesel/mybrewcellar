package pages

import (
	"models"
	//"github.com/skiesel/mybrewcellar/models"
	"appengine"
	"appengine/user"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type responseJson struct {
	Status string
	Error  string
	Data   interface{}
}

func init() {
	http.HandleFunc("/api/update-account", handleUpdateAccount)

	http.HandleFunc("/api/new-cellar", handleNewCellarRequest)
	http.HandleFunc("/api/delete-cellar", handleDeleteCellarRequest)

	http.HandleFunc("/api/new-beer", handleNewBeerRequest)
	http.HandleFunc("/api/delete-beer", handleDeleteBeerRequest)
	http.HandleFunc("/api/transfer-beer", handleTransferBeerRequest)

	http.HandleFunc("/api/new-tasting", handleNewTastingRequest)
	http.HandleFunc("/api/delete-tasting", handleDeleteTastingRequest)
}

func getAccount(c appengine.Context) *models.Account {
	usr := user.Current(c)
	if usr != nil {
		return models.GetAccount(c, usr.Email)
	}
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	obj := responseJson{
		Status: "FAILURE",
		Error:  err.Error(),
		Data:   "",
	}

	response, err := json.Marshal(obj)

	if err != nil {
		panic(err)
	}

	w.Write(response)
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	obj := responseJson{
		Status: "SUCCESS",
		Error:  "",
		Data:   data,
	}

	response, err := json.Marshal(obj)

	if err != nil {
		writeError(w, err)
	} else {
		w.Write(response)
	}
}

type simpleUser struct {
	Username string
}

func handleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	newUsername := r.PostFormValue("username")

	if newUsername != "" {
		account := getAccount(c)
		if account == nil {
			writeError(w, errors.New("no account"))
			return
		}

		account.User.UserID = newUsername
		err := models.SaveAccount(c, account)
		if err != nil {
			writeError(w, err)
			return
		}
		writeSuccess(w, simpleUser{Username: newUsername})
		return
	}

	writeError(w, errors.New("No username supplied"))
}

type simpleCellar struct {
	ID   int
	Name string
	BeerCount int
}

func handleNewCellarRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	newCellar := r.PostFormValue("cellarName")

	if newCellar != "" {
		account := getAccount(c)
		if account == nil {
			writeError(w, errors.New("no account"))
			return
		}

		err := account.AddCellar(newCellar)
		if err != nil {
			writeError(w, err)
			return
		}

		err = models.SaveAccount(c, account)
		if err != nil {
			writeError(w, err)
			return
		}

		cellar := account.Cellars[newCellar]
		writeSuccess(w, simpleCellar{ID: cellar.ID, Name: cellar.Name, BeerCount : 0})
		return
	}

	writeError(w, errors.New("No cellar supplied"))
}

func handleDeleteCellarRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	deleteCellarID, err := strconv.Atoi(r.PostFormValue("cellarID"))
	if err != nil {
		writeError(w, err)
		return
	}

	account := getAccount(c)
	if account == nil {
		writeError(w, errors.New("no account"))
	}

	cellar := account.CellarsByID[deleteCellarID]
	err = account.DeleteCellarByID(deleteCellarID)
	if err != nil {
		writeError(w, err)
		return
	}

	err = models.SaveAccount(c, account)
	if err != nil {
		writeError(w, err)
		return
	}

	writeSuccess(w, simpleCellar{ID: cellar.ID, Name: cellar.Name})
}

type simpleBeer struct {
	ID       int
	Name     string
	AverageRating int
	Quantity int
	Notes    string
	Brewed   string
	Added    string
	Age      string
}

func handleNewBeerRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	cellarID, err := strconv.Atoi(r.PostFormValue("cellarID"))
	if err != nil {
		writeError(w, err)
		return
	}

	quantity, err := strconv.Atoi(r.PostFormValue("quantity"))
	if err != nil {
		writeError(w, err)
		return
	}

	name := r.PostFormValue("name")
	if name == "" {
		writeError(w, errors.New("No beer supplied"))
		return
	}

	notes := r.PostFormValue("notes")

	brewed := r.PostFormValue("brewed")
	var brewDate *models.Date
	if brewed == "" {
		brewDate = models.Now()
	} else {
		brewDate = models.ParseDate(brewed)
	}

	added := r.PostFormValue("added")
	var addedDate *models.Date
	if added == "" {
		addedDate = models.Now()
	} else {
		addedDate = models.ParseDate(added)
	}

	account := getAccount(c)
	if account == nil {
		writeError(w, errors.New("no account"))
		return
	}

	cellar := account.CellarsByID[cellarID]
	if cellar == nil {
		writeError(w, errors.New("cellar not found"))
		return
	}

	beer := &models.Beer{
		ID:            cellar.NextBeerID,
		Name:          name,
		Notes:         notes,
		Brewed:        brewDate,
		Added:         addedDate,
		Quantity:      quantity,
		NextTastingID: 0,
		TastingsByID:  map[int]*models.Tasting{},
	}

	cellar.NextBeerID++

	cellar.Beers[beer.Name] = beer
	cellar.BeersByID[beer.ID] = beer

	err = models.SaveAccount(c, account)
	if err != nil {
		writeError(w, err)
		return
	}

	writeSuccess(w, simpleBeer{
		ID:       beer.ID,
		Name:     beer.Name,
		AverageRating: 0,
		Quantity: quantity,
		Notes:    beer.Notes,
		Brewed:   beer.Brewed.ToString(),
		Added:    beer.Added.ToString(),
		Age:      beer.GetAgeString(),
	})
}

func handleDeleteBeerRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	cellarID, err := strconv.Atoi(r.PostFormValue("cellarID"))
	if err != nil {
		writeError(w, err)
		return
	}

	beerID, err := strconv.Atoi(r.PostFormValue("beerID"))
	if err != nil {
		writeError(w, err)
		return
	}

	account := getAccount(c)
	if account == nil {
		writeError(w, errors.New("no account"))
		return
	}

	cellar := account.CellarsByID[cellarID]
	if cellar == nil {
		writeError(w, errors.New("cellar not found"))
		return
	}

	beer := cellar.BeersByID[beerID]
	if beer == nil {
		writeError(w, errors.New("beer not found"))
		return
	}

	delete(cellar.Beers, beer.Name)
	delete(cellar.BeersByID, beer.ID)

	err = models.SaveAccount(c, account)
	if err != nil {
		writeError(w, err)
		return
	}

	writeSuccess(w, simpleBeer{ID: beer.ID, Name: beer.Name})
}

type simpleTransfer struct {
	FromCellar simpleCellar
	ToCellar   simpleCellar
	Beer       simpleBeer
}

func handleTransferBeerRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	fromCellarID, err := strconv.Atoi(r.PostFormValue("fromCellarID"))
	if err != nil {
		writeError(w, err)
		return
	}

	toCellarID, err := strconv.Atoi(r.PostFormValue("toCellarID"))
	if err != nil {
		writeError(w, err)
		return
	}

	beerID, err := strconv.Atoi(r.PostFormValue("beerID"))
	if err != nil {
		writeError(w, err)
		return
	}

	account := getAccount(c)
	if account == nil {
		writeError(w, errors.New("no account"))
		return
	}

	fromCellar := account.CellarsByID[fromCellarID]
	if fromCellar == nil {
		writeError(w, errors.New("Could not move from cellar, cellar not found"))
		return
	}
	toCellar := account.CellarsByID[toCellarID]
	if toCellar == nil {
		writeError(w, errors.New("Could not move to cellar, cellar not found"))
		return
	}
	beer := fromCellar.BeersByID[beerID]
	if beer == nil {
		writeError(w, errors.New("Could not move from cellar, beer not found"))
		return
	}

	delete(fromCellar.Beers, beer.Name)
	delete(fromCellar.BeersByID, beer.ID)

	oldBeerID := beer.ID

	beer.ID = toCellar.NextBeerID
	toCellar.NextBeerID++

	toCellar.Beers[beer.Name] = beer
	toCellar.BeersByID[beer.ID] = beer

	err = models.SaveAccount(c, account)
	if err != nil {
		writeError(w, err)
		return
	}
	writeSuccess(w, simpleTransfer{
		FromCellar: simpleCellar{
			ID:   fromCellar.ID,
			Name: fromCellar.Name,
		},
		ToCellar: simpleCellar{
			ID:   toCellar.ID,
			Name: toCellar.Name,
		},
		Beer: simpleBeer{
			ID:   oldBeerID,
			Name: beer.Name,
		},
	})
}

type simpleTasting struct {
	ID int
	Rating int
	Notes string
	TastedDate string
	AgeTastedDate string
	Quantity int
	AverageRating float64
}

func handleNewTastingRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	cellarID, err := strconv.Atoi(r.PostFormValue("cellarID"))
	if err != nil {
		writeError(w, err)
		return
	}

	beerID, err := strconv.Atoi(r.PostFormValue("beerID"))
	if err != nil {
		writeError(w, err)
		return
	}

	tasted := r.PostFormValue("tasted")
	var tastedDate *models.Date
	if tasted == "" {
		tastedDate = models.Now()
	} else {
		tastedDate = models.ParseDate(tasted)
	}

	rating, err := strconv.Atoi(r.PostFormValue("rating"))
	if err != nil {
		writeError(w, err)
		return
	}

	decrement := r.PostFormValue("decrement") == "yes"
	notes := r.PostFormValue("notes")

	account := getAccount(c)
	if account == nil {
		writeError(w, errors.New("no account"))
		return
	}

	cellar := account.CellarsByID[cellarID]
	if cellar == nil {
		writeError(w, errors.New("cellar not found"))
		return
	}

	beer := cellar.BeersByID[beerID]
	if beer == nil {
		writeError(w, errors.New("beer not found"))
		return
	}	

	tasting := &models.Tasting {
		ID : beer.NextTastingID,
		Rating : rating,
		Notes : notes,
		Date : tastedDate,
	}
	
	if decrement {
		beer.Quantity--
		if beer.Quantity < 0 {
			beer.Quantity = 0
		}
	}
	beer.NextTastingID++

	beer.TastingsByID[tasting.ID] = tasting

	err = models.SaveAccount(c, account)
	if err != nil {
		writeError(w, err)
		return
	}
	
	writeSuccess(w, simpleTasting{
		ID : tasting.ID,
		Rating : tasting.Rating,
		Notes : tasting.Notes,
		TastedDate : tasting.Date.ToString(),
		AgeTastedDate : beer.GetTastingAge(tasting),
		Quantity : beer.Quantity,
		AverageRating : beer.GetAverageRating(),
	})
}

func handleDeleteTastingRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	cellarID, err := strconv.Atoi(r.PostFormValue("cellarID"))
	if err != nil {
		writeError(w, err)
		return
	}

	beerID, err := strconv.Atoi(r.PostFormValue("beerID"))
	if err != nil {
		writeError(w, err)
		return
	}

	tastingID, err := strconv.Atoi(r.PostFormValue("tastingID"))
	if err != nil {
		writeError(w, err)
		return
	}

	account := getAccount(c)
	if account == nil {
		writeError(w, errors.New("no account"))
		return
	}

	cellar := account.CellarsByID[cellarID]
	if cellar == nil {
		writeError(w, errors.New("cellar not found"))
		return
	}

	beer := cellar.BeersByID[beerID]
	if beer == nil {
		writeError(w, errors.New("beer not found"))
		return
	}

	tasting := beer.TastingsByID[tastingID]
	if tasting == nil {
		writeError(w, errors.New("tasting not found"))
		return
	}

	delete(beer.TastingsByID, tasting.ID)

	err = models.SaveAccount(c, account)
	if err != nil {
		writeError(w, err)
		return
	}	

	writeSuccess(w, simpleTasting{ ID : tasting.ID,
		Quantity : beer.Quantity,
		AverageRating : beer.GetAverageRating(),
	})
}