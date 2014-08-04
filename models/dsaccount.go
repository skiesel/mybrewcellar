package models

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

func SaveAccount(c appengine.Context, account *Account) error {

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

		cellarQuery := datastore.NewQuery("Cellar").Ancestor(accountKey)
		cellarResults := cellarQuery.Run(c)
		oldCellarKeys := map[int]*datastore.Key{}
		for {
			cellarDS := &CellarDS{}
			cellarKey, err := cellarResults.Next(cellarDS)
			if err == datastore.Done {
				c.Errorf("no more cellars")
				break
			}
			if err != nil {
				c.Errorf("fetching next Cellar: %v", err)
				break
			}
			oldCellarKeys[cellarDS.ID] = cellarKey
		}

		cellarKeys := make([]*datastore.Key, len(account.Cellars))
		cellarDSs := make([]*CellarDS, len(account.Cellars))
		beerCount := 0
		tastingCount := 0

		for cellarID, datastoreKey := range oldCellarKeys {
			if account.CellarsByID[cellarID] == nil {
				datastore.Delete(c, datastoreKey)
			}
		}

		i := 0
		for _, cellar := range account.Cellars {
			existingKey := oldCellarKeys[cellar.ID]
			if existingKey != nil {
				cellarKeys[i] = existingKey
			} else {
				cellarKeys[i] = datastore.NewIncompleteKey(c, "Cellar", accountKey)
			}
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

			oldBeerKeys := map[int]*datastore.Key{}
			beerQuery := datastore.NewQuery("Beer").Ancestor(cellarKeys[i])
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
				oldBeerKeys[beerDS.ID] = beerKey
			}

			for beerID, datastoreKey := range oldBeerKeys {
				if cellar.BeersByID[beerID] == nil {
					datastore.Delete(c, datastoreKey)
				}
			}

			for _, beer := range cellar.Beers {
				existingKey, exists := oldBeerKeys[beer.ID]
				if exists {
					beerKeys[curBeerIndex] = existingKey
				} else {
					beerKeys[curBeerIndex] = datastore.NewIncompleteKey(c, "Beer", cellarKeys[i])
				}
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

				oldTastingKeys := map[int]*datastore.Key{}
				tastingQuery := datastore.NewQuery("Tasting").Ancestor(beerKeys[curBeerIndex])
				tastingResults := tastingQuery.Run(c)
				for {
					tastingDS := &TastingDS{}
					tastingKey, err := tastingResults.Next(tastingDS)
					if err == datastore.Done {
						break
					}
					if err != nil {
						c.Errorf("fetching next Tasting: %v", err)
						break
					}
					oldTastingKeys[tastingDS.ID] = tastingKey
				}

				for tastingID, datastoreKey := range oldTastingKeys {
					if beer.TastingsByID[tastingID] == nil {
						datastore.Delete(c, datastoreKey)
					}
				}

				for _, tasting := range beer.TastingsByID {
					existingKey, exists := oldTastingKeys[tasting.ID]
					if exists {
						tastingKeys[curTastingIndex] = existingKey
					} else {
						tastingKeys[curTastingIndex] = datastore.NewIncompleteKey(c, "Tasting", beerKeys[curBeerIndex])
					}
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

		c.Errorf("%s", cellar.Name)

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

func DeleteAccount(c appengine.Context, account *Account) error {
	accountKey := datastore.NewKey(c, "Account", account.User.Email, 0, nil)
	err := datastore.Delete(c, accountKey)
	return err
}
