package db

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/olivere/elastic"
)

// ECObject represents a remove ElasticSearch object
type ECObject struct {
	client *elastic.Client

	index string
	tp    string
	id    string

	Fields map[string]interface{}
}

// IsIndexed checks if an object is still indexed
func (obj *ECObject) IsIndexed() bool {
	return obj.id != ""
}

// Index indexes this object in the database
func (obj *ECObject) Index() (err error) {
	// if we already have an id, it is already indexed
	if obj.IsIndexed() {
		return errors.New("Already indexed")
	}

	// perform the Indexing operation
	ctx := context.Background()
	result, err := obj.client.Index().Index(obj.index).Type(obj.tp).BodyJson(obj.Fields).Do(ctx)
	if err == nil && result.Shards.Successful <= 0 {
		err = errors.New("Index() reported 0 successful shards")
	}

	if err != nil {
		return
	}

	// grab the new object id
	obj.id = result.Id

	return
}

// Reload reloads this object from the database
func (obj *ECObject) Reload() (err error) {

	if !obj.IsIndexed() {
		return errors.New("Not indexed")
	}

	ctx := context.Background()

	// fetch from the db and return unless an error occured
	result, err := obj.client.Get().Index(obj.index).Type(obj.tp).Id(obj.id).Do(ctx)
	if err == nil && !result.Found {
		err = errors.New("Reload() did not find object")
	}

	if err != nil {
		err = obj.setSource(result.Source)
	}

	return
}

func (obj *ECObject) setSource(source *json.RawMessage) error {
	return json.Unmarshal(*source, &obj.Fields)
}

// Save saves this object into the database
func (obj *ECObject) Save() (err error) {
	if !obj.IsIndexed() {
		return errors.New("Not indexed")
	}

	// replace the entire item in the database
	ctx := context.Background()
	res, err := obj.client.Update().Index(obj.index).Type(obj.tp).Id(obj.id).Doc(obj.Fields).Do(ctx)

	if err == nil && (res.Result != "noop" && res.Shards.Successful <= 0) {
		err = errors.New("Save() reported non-noop result with 0 successful shards ")
	}

	return
}

// Delete deletes this object
func (obj *ECObject) Delete() (err error) {
	if !obj.IsIndexed() {
		return errors.New("Not indexed")
	}

	// just clears the object
	ctx := context.Background()
	res, err := obj.client.Delete().Index(obj.index).Type(obj.tp).Id(obj.id).Do(ctx)

	if err == nil && res.Result != "deleted" {
		err = errors.New("Delete() did not report deleted result ")
	}

	// id no longer valid
	if err == nil {
		obj.id = ""
	}

	return
}
