package db

import (
	"time"

	"github.com/olivere/elastic"
)

// Connect connects to the server
func Connect(retryInterval time.Duration, onRetry func(error), funcs ...elastic.ClientOptionFunc) (cli *elastic.Client) {
	var err error
	for {
		cli, err = elastic.NewClient(funcs...)

		if err == nil {
			break
		}

		// and wait for next time
		onRetry(err)
		time.Sleep(retryInterval)
	}
	return
}
