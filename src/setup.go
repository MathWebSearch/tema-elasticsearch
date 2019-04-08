package src

import "fmt"

// SetupDB sets up the database or errors out
func (ec *ElasticClient) SetupDB() (err error) {
	// check if the index exists
	exists, err := ec.client.IndexExists(ec.index).Do(*ec.ctx)
	if err != nil {
		return
	}

	// if yes, do nothing
	if exists {
		fmt.Printf("Index %q exists, skipping setup\n", ec.index)

		// if no, create it
	} else {
		fmt.Printf("Creating index %q \n", ec.index)
		_, err := ec.client.CreateIndex(ec.index).BodyString(indexMapping).Do(*ec.ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

const indexMapping = `
{
    "settings": {
        "index": {
            "refresh_interval": "-1",
            "number_of_replicas": 0,
            "number_of_shards": 128
        }
    },
    "mappings": {
        "doc": {
            "dynamic": false,
            "properties": {
                "metadata": {
                    "dynamic": true,
                    "type": "object"
                },
                "mws_ids": {
                    "type": "long",
                    "store": false
                },
                "text": {
                    "type": "text",
                    "store": false
                },
                "mws_id": {
                    "enabled": false,
                    "type": "object"
                },
                "math": {
                    "enabled": false,
                    "type": "object"
                }
            }
        }
    }
}`
