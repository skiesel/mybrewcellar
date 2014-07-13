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
		writeSuccess(w, simpleCellar{ID: cellar.ID, Name: cellar.Name})
		return
	}

	writeError(w, errors.New("No cellar supplied"))
}

func handleDeleteCellarRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	newCellar := r.PostFormValue("cellarName")

	if newCellar != "" {
		account := getAccount(c)
		if account == nil {
			writeError(w, errors.New("no account"))
		}

		cellar := account.Cellars[newCellar]
		err := account.DeleteCellar(newCellar)
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
		return
	}

	writeError(w, errors.New("No cellar supplied"))
}

type simpleBeer struct {
	ID int
	Name string
	Notes string
	Brewed string
	Added string
	Quantity int
}

func handleNewBeerRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("add beer")
}

func handleDeleteBeerRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	cellarId, err := strconv.Atoi(r.PostFormValue("cellarID"))
	if err != nil {
		writeError(w, err)
		return
	}
	
	beerName := r.PostFormValue("beerName")
	
	account := getAccount(c)
	if account == nil {
		writeError(w, errors.New("no account"))
		return
	}

	cellar := account.CellarsByID[cellarId]
	if cellar == nil {
		writeError(w, errors.New("cellar not found"))
		return	
	}

	beer := cellar.Beers[beerName]
	if beer == nil {
		writeError(w, errors.New("beer not found"))
		return	
	}

	for key, cellar := range account.Cellars {
		c.Infof("%s : %s", key, cellar.Name)
		for key, beer := range cellar.Beers {
			c.Infof("\t%s : %s", key, beer.Name)
		}
		for id, beer := range cellar.BeersByID {
			c.Infof("\t%d : %s", id, beer.Name)
		}
	}
	for id, cellar := range account.CellarsByID {
		c.Infof("%d : %s", id, cellar.Name)
		for key, beer := range cellar.Beers {
			c.Infof("\t%s : %s", key, beer.Name)
		}
		for id, beer := range cellar.BeersByID {
			c.Infof("\t%d : %s", id, beer.Name)
		}
	}

	if account.CellarsByID[cellarId] != account.Cellars[account.CellarsByID[cellarId].Name] {
		panic("NOT THE SAME CELLAR")
	}

	delete(cellar.Beers, beer.Name)
	delete(cellar.BeersByID, beer.ID)

	c.Infof("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	for key, cellar := range account.Cellars {
		c.Infof("%s : %s", key, cellar.Name)
		for key, beer := range cellar.Beers {
			c.Infof("\t%s : %s", key, beer.Name)
		}
		for id, beer := range cellar.BeersByID {
			c.Infof("\t%d : %s", id, beer.Name)
		}
	}
	for id, cellar := range account.CellarsByID {
		c.Infof("%d : %s", id, cellar.Name)
		for key, beer := range cellar.Beers {
			c.Infof("\t%s : %s", key, beer.Name)
		}
		for id, beer := range cellar.BeersByID {
			c.Infof("\t%d : %s", id, beer.Name)
		}
	}

	err = models.SaveAccount(c, account)
	if err != nil {
		writeError(w, err)
		return
	}

	writeSuccess(w, simpleBeer{ID: beer.ID, Name: beer.Name})	
}

