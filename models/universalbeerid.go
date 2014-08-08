package models

import (
	"appengine"
	"appengine/datastore"
)

type Counter struct {
	Count int
}

func getKey(c appengine.Context) *datastore.Key {
  return datastore.NewKey(c, "ShardHead", "ShardHead", 0, nil)
}

func GetAndIncrementUniversalBeerID(c appengine.Context) (int, error) {
	id := -1
	err := datastore.RunInTransaction(c, func(c appengine.Context) error {
    key := getKey(c)
    counter := &Counter{}
    err := datastore.Get(c, key, counter)

    if err != nil && err != datastore.ErrNoSuchEntity {
      return err
    }

    id = counter.Count
    counter.Count++

    _, err = datastore.Put(c, key, counter)

    return err
	}, nil)
	if err != nil {
		return -1, err
	}

	return id, err
}
