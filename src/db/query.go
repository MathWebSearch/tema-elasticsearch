package db

import (
	"context"
	"errors"

	"github.com/olivere/elastic"
)

// DeleteBulk deletes objects subject to a given query
func DeleteBulk(client *elastic.Client, index string, tp string, query elastic.Query) (err error) {
	ctx := context.Background()
	res, err := client.DeleteByQuery(index).Type(tp).Query(query).Do(ctx)

	if err == nil && res.TimedOut {
		err = errors.New("DeleteByQuery() reported TimedOut=true")
	}

	return
}

// UpdateAll updates all objects inside a given index
func UpdateAll(client *elastic.Client, index string, tp string, script *elastic.Script) (err error) {
	ctx := context.Background()
	res, err := client.UpdateByQuery(index).Type(tp).Query(elastic.NewMatchAllQuery()).Script(script).Do(ctx)

	if err == nil && res.TimedOut {
		err = errors.New("UpdateByQuery() reported TimedOut=true")
	}

	return
}

// CreateBulk creates a lazily computed set of objects in bulk
func CreateBulk(client *elastic.Client, index string, tp string, objects chan map[string]interface{}) (err error) {
	// create a new bulk request
	bulkRequest := client.Bulk()

	for object := range objects {
		req := elastic.NewBulkIndexRequest().Index(index).Type(tp).Doc(object)
		bulkRequest.Add(req)
	}

	ctx := context.Background()
	res, err := bulkRequest.Do(ctx)

	if err == nil && res.Errors {
		err = errors.New("Bulk() reported Errors=true")
	}

	return
}

// CreateIndex creates an index unless it already exists
func CreateIndex(client *elastic.Client, index string, mapping interface{}) (created bool, err error) {
	ctx := context.Background()

	// check if the index exists
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return
	}

	// create it if not
	if !exists {
		res, err := client.CreateIndex(index).BodyJson(mapping).Do(ctx)
		if err == nil && !res.Acknowledged {
			err = errors.New("CreateIndex() reported acknowledged=false")
		}

		if err != nil {
			return false, err
		}
		created = true
	}

	return
}
