package src

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic"
	"github.com/olivere/elastic/config"
)

// ElasticClient represents a wrapped elastic client
type ElasticClient struct {
	client *elastic.Client
	ctx    *context.Context

	index string
	url   string
	dir   string
}

// MakeClientAndWait makes a client and waits for it
func MakeClientAndWait(args *Args) (client *ElasticClient) {
	// build the url
	url := args.ElasticURL()

	var cli *elastic.Client
	var err error

	fmt.Println("Waiting for Elasticsarch to start up ...")
	for {
		cli, err = elastic.NewClientFromConfig(&config.Config{URL: url})

		if err == nil {
			break
		}

		fmt.Printf("Failed to connect to %q, trying again is 5 second(s)\n", url)

		// and wait for next time
		time.Sleep(5 * time.Second)
	}
	fmt.Printf("Connected successfully to %q\n", url)

	ctx := context.Background()
	client = &ElasticClient{cli, &ctx, args.ElasticIndex, url, args.IndexDir}

	return
}
