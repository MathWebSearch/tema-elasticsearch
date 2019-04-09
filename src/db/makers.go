package db

import (
	"context"
	"errors"

	"github.com/olivere/elastic"
)

// GetObject fetches a new EC object from the server
func GetObject(client *elastic.Client, index string, tp string, id string) (obj *ECObject, err error) {
	// create an empty object
	obj = &ECObject{client, index, tp, id, nil}

	// reload it from the db, clear if it fails
	err = obj.Reload()
	if err != nil {
		obj = nil
	}
	return
}

// CreateObject creates a new ec object on the server
func CreateObject(client *elastic.Client, index string, tp string, Fields map[string]interface{}) (obj *ECObject, err error) {
	obj = &ECObject{client, index, tp, "", Fields}

	err = obj.Index()
	if err != nil {
		obj = nil
	}

	return
}

// FetchObjects fetches objects subject to an exact query
func FetchObjects(client *elastic.Client, index string, tp string, query elastic.Query) (<-chan *ECObject, error) {
	// run the query in the background
	ctx := context.Background()
	result, err := client.Search(index).Type(tp).Query(query).Do(ctx)

	if err == nil && result.TimedOut {
		err = errors.New("Search() returned TimedOut=true")
	}

	// bail out if we have an error
	if err != nil {
		return nil, err
	}

	results := make(chan *ECObject)

	go func() {
		for _, hit := range result.Hits.Hits {
			obj := &ECObject{client, index, tp, hit.Id, nil}
			obj.setSource(hit.Source)
			results <- obj
		}
		close(results)
	}()

	return results, nil
}

// FetchObject fetches a single object from the database or returns nil
func FetchObject(client *elastic.Client, index string, tp string, query elastic.Query) (obj *ECObject, err error) {
	// make a query
	results, err := FetchObjects(client, index, tp, query)
	if err != nil {
		return
	}

	// fetch the candidate
	for candidate := range results {
		obj = candidate
		break
	}

	// empty the channel
	go func() {
		for range results {
		}
	}()

	// and return the result (if any)
	return

}

// FetchOrCreateObject fetches the object returned from the query, or creates a new one if no result is retrieved
func FetchOrCreateObject(client *elastic.Client, index string, tp string, query elastic.Query, NewFields map[string]interface{}) (obj *ECObject, created bool, err error) {
	// first try and fetch the object
	obj, err = FetchObject(client, index, tp, query)
	if err != nil || obj != nil {
		return
	}

	// if that fails create it
	obj, err = CreateObject(client, index, tp, NewFields)
	if err != nil {
		created = true
	}

	return
}
