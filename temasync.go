package main

import (
	"fmt"
	"os"
	"time"

	"github.com/MathWebSearch/tema-elasticsync/src"
	"github.com/MathWebSearch/tema-elasticsync/src/db"
	"github.com/MathWebSearch/tema-elasticsync/src/sync"
	"github.com/olivere/elastic"
)

func main() {
	// parse and validate arguments
	args := src.ParseArgs(os.Args)
	if !args.Validate() {
		os.Exit(1)
		return
	}

	// connect to the database
	url := args.ElasticURL()
	fmt.Printf("Connecting to %q ...\n", url)

	client := db.Connect(5*time.Second, func(e error) {
		fmt.Printf("  Unable to connect: %s. Trying again in 5 seconds. \n", e)
	}, elastic.SetURL(url), elastic.SetSniff(false))
	fmt.Println("Connected. ")

	// make a sync process
	process := sync.NewProcess(client, args.IndexDir)
	process.Run()

	fmt.Println("Finished, exiting. ")
}

func die(err error) {

	if err != nil {
		panic(err)
	} else {
		panic("Something went wrong")
	}
}

// TODO: 1. Wait for tema-search to be up on the given port
// 2. Check if we have to run setup
// 3. Hash the directory; if it has changed clear out and fully re-index
