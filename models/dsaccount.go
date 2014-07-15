package models

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

type AccountDS struct {
	UserID       string
	UserEmail    string
	NextCellarID int
}

type CellarDS struct {
	ID         int
	NextBeerID int
	Name       string
}

type BeerDS struct {
	ID            int
	Name          string
	Notes         string
	Brewed        string
	Added         string
	Quantity      int
	NextTastingID int
}

type TastingDS struct {
	ID     int
	Rating int
	Notes  string
	Date   string
}

func SaveAccount(c appengine.Context, account *Account) error {

	for _, cellar := range account.CellarsByID {
		if cellar != account.Cellars[cellar.Name] {
			panic("Save Account: NOT THE SAME CELLAR")
		}
	}

	cacheItem := &memcache.Item{
		Key:    account.User.Email,
		Object: account,
	}

	err := memcache.Gob.Set(c, cacheItem)
	if err != nil {
		c.Infof("Account not cached (%v)", err)
	}

	err = datastore.RunInTransaction(c, func(c appengine.Context) error {

		accountKey := datastore.NewKey(c, "Account", account.User.Email, 0, nil)
		accountDS := account.toAccountDS()

		_, err := datastore.Put(c, accountKey, accountDS)
		if err != nil {
			c.Errorf("1 %s", err.Error())
			return err
		}

		if len(account.Cellars) <= 0 {
			return nil
		}

		cellarKeys := make([]*datastore.Key, len(account.Cellars))
		cellarDSs := make([]*CellarDS, len(account.Cellars))
		beerCount := 0
		tastingCount := 0

		i := 0
		for _, cellar := range account.Cellars {
			cellarKeys[i] = datastore.NewIncompleteKey(c, "Cellar", accountKey)
			cellarDSs[i] = cellar.toCellarDS()
			i++
			beerCount += len(cellar.Beers)
			for _, beer := range cellar.Beers {
				tastingCount += len(beer.TastingsByID)
			}
		}

		cellarKeys, err = datastore.PutMulti(c, cellarKeys, cellarDSs)
		if err != nil {
			c.Errorf("2 %s", err.Error())
			return err
		}

		beerKeys := make([]*datastore.Key, beerCount)
		beerDSs := make([]*BeerDS, beerCount)

		tastingKeys := make([]*datastore.Key, tastingCount)
		tastingDSs := make([]*TastingDS, tastingCount)

		curBeerIndex := 0

		for i, cellarDS := range cellarDSs {
			cellar := account.CellarsByID[cellarDS.ID]

			for _, beer := range cellar.Beers {
				beerKeys[curBeerIndex] = datastore.NewIncompleteKey(c, "Beer", cellarKeys[i])
				beerDSs[curBeerIndex] = beer.toBeerDS()
				curBeerIndex++
			}
		}

		beerKeys, err = datastore.PutMulti(c, beerKeys, beerDSs)
		if err != nil {
			c.Errorf("3 %s", err.Error())
			return err
		}

		curBeerIndex = 0
		curTastingIndex := 0

		for _, cellarDS := range cellarDSs {
			cellar := account.CellarsByID[cellarDS.ID]

			for _, beer := range cellar.Beers {

				for _, tasting := range beer.TastingsByID {
					tastingKeys[curTastingIndex] = datastore.NewIncompleteKey(c, "Tasting", beerKeys[curBeerIndex])
					tastingDSs[curTastingIndex] = tasting.toTastingDS()
					curTastingIndex++
				}

				curBeerIndex++
			}
		}

		if tastingCount > 0 {
			_, err = datastore.PutMulti(c, tastingKeys, tastingDSs)
		}

		return err
	}, nil)

	return err
}

func (account Account) toAccountDS() *AccountDS {
	return &AccountDS{
		UserID:       account.User.UserID,
		UserEmail:    account.User.Email,
		NextCellarID: account.NextCellarID,
	}
}

func (cellar Cellar) toCellarDS() *CellarDS {
	return &CellarDS{
		ID:         cellar.ID,
		NextBeerID: cellar.NextBeerID,
		Name:       cellar.Name,
	}
}

func (beer *Beer) toBeerDS() *BeerDS {
	return &BeerDS{
		ID:            beer.ID,
		Name:          beer.Name,
		Notes:         beer.Notes,
		Brewed:        beer.Brewed.ToDSString(),
		Added:         beer.Added.ToDSString(),
		Quantity:      beer.Quantity,
		NextTastingID: beer.NextTastingID,
	}
}

func (tasting *Tasting) toTastingDS() *TastingDS {
	return &TastingDS{
		ID:     tasting.ID,
		Rating: tasting.Rating,
		Notes:  tasting.Notes,
		Date:   tasting.Date.ToDSString(),
	}
}

func GetAccount(c appengine.Context, email string) *Account {

	cachedAccount := &Account{}
	_, err := memcache.Gob.Get(c, email, cachedAccount)
	if err == nil {
		c.Infof("Cache Hit!")

		// avoid problems with memcache doing a deep copy and not
		// maintaining references properly
		for _, cellar := range cachedAccount.CellarsByID {
			cachedAccount.Cellars[cellar.Name] = cellar
			for _, beer := range cellar.BeersByID {
				cellar.Beers[beer.Name] = beer
			}
		}

		return cachedAccount
	} else {
		c.Infof("Cache Miss, retrieving from datastore! (%v)", err)
	}

	c.Infof("~~~~~~~LoadAccount~~~~~~~")

	accountKey := datastore.NewKey(c, "Account", email, 0, nil)
	accountDS := &AccountDS{}
	err = datastore.Get(c, accountKey, accountDS)
	if err != nil {
		c.Errorf("Did not find account: %v (%v)", email, err.Error())
		return nil
	}

	account := accountDS.toAccount()

	cellarQuery := datastore.NewQuery("Cellar").Ancestor(accountKey)
	cellarResults := cellarQuery.Run(c)
	for {
		cellarDS := &CellarDS{}
		cellarKey, err := cellarResults.Next(cellarDS)
		if err == datastore.Done {
			break
		}
		if err != nil {
			c.Errorf("fetching next Cellar: %v", err)
			break
		}

		cellar := cellarDS.toCellar()

		account.Cellars[cellarDS.Name] = cellar
		account.CellarsByID[cellar.ID] = cellar

		beerQuery := datastore.NewQuery("Beer").Ancestor(cellarKey)
		beerResults := beerQuery.Run(c)
		for {
			beerDS := &BeerDS{}
			beerKey, err := beerResults.Next(beerDS)
			if err == datastore.Done {
				break
			}
			if err != nil {
				c.Errorf("fetching next Beer: %v", err)
				break
			}

			beer := beerDS.toBeer()

			cellar.Beers[beer.Name] = beer
			cellar.BeersByID[beer.ID] = beer

			tastingQuery := datastore.NewQuery("Tasting").Ancestor(beerKey)
			tastingResults := tastingQuery.Run(c)
			for {
				tastingDS := &TastingDS{}
				_, err := tastingResults.Next(tastingDS)
				if err == datastore.Done {
					break
				}
				if err != nil {
					c.Errorf("fetching next Tasting: %v", err)
					break
				}

				tasting := tastingDS.toTasting()
				beer.TastingsByID[tasting.ID] = tasting
			}
		}
	}

	return account
}

func (accountDS *AccountDS) toAccount() *Account {
	return &Account{
		User: &User{
			UserID: accountDS.UserID,
			Email:  accountDS.UserEmail,
		},
		NextCellarID: accountDS.NextCellarID,
		Cellars:      map[string]*Cellar{},
		CellarsByID:  map[int]*Cellar{},
	}
}

func (cellarDS *CellarDS) toCellar() *Cellar {
	return &Cellar{
		ID:         cellarDS.ID,
		NextBeerID: cellarDS.NextBeerID,
		Name:       cellarDS.Name,
		Beers:      map[string]*Beer{},
		BeersByID:  map[int]*Beer{},
	}
}

func (beerDS *BeerDS) toBeer() *Beer {
	return &Beer{
		ID:            beerDS.ID,
		Name:          beerDS.Name,
		Notes:         beerDS.Notes,
		Brewed:        ParseDate(beerDS.Brewed),
		Added:         ParseDate(beerDS.Added),
		Quantity:      beerDS.Quantity,
		NextTastingID: beerDS.NextTastingID,
		TastingsByID:  map[int]*Tasting{},
	}
}

func (tastingDS *TastingDS) toTasting() *Tasting {
	return &Tasting{
		ID:     tastingDS.ID,
		Rating: tastingDS.Rating,
		Notes:  tastingDS.Notes,
		Date:   ParseDate(tastingDS.Date),
	}
}

func DeleteAccount(c appengine.Context, account *Account) error {
	accountKey := datastore.NewKey(c, "Account", account.User.Email, 0, nil)
	err := datastore.Delete(c, accountKey)
	return err
}
