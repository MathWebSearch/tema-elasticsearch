package src

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/olivere/elastic"
)

// LoadHarvests loads the hervests from the appropriate directory
func (ec *ElasticClient) LoadHarvests() (err error) {
	// for now we just bulk load all the existing files
	return iterateFiles(ec.dir, ".json", ec.loadBulkFile)
}

type temaIndexMeta struct {
	Index struct {
		ID string `json:"_id"`
	} `json:"index"`
}

func (ec *ElasticClient) loadBulkFile(filename string) (err error) {
	// load a new harvest
	fmt.Printf("Loading harvest from %q\n", filename)

	// create a new bulk request
	bulkRequest := ec.client.Bulk()

	err = processLinePairs(filename, true, func(meta string, doc string) error {
		// read the meta index
		var metaIndex temaIndexMeta
		err := json.Unmarshal([]byte(meta), &metaIndex)
		if err != nil {
			return err
		}

		// read the document
		var docParse interface{}
		err = json.Unmarshal([]byte(doc), &docParse)
		if err != nil {
			return err
		}

		// create a new request
		req := elastic.NewBulkIndexRequest().Index(ec.index).Type("doc").Id(metaIndex.Index.ID).Doc(docParse)
		bulkRequest = bulkRequest.Add(req)

		return nil
	})
	if err != nil {
		return
	}

	// and do the bulk request
	_, err = bulkRequest.Do(*ec.ctx)
	return
}

type dualParser func(string, string) error

const bufferCapcityInBytes = 128 * 1024 // 128 MB

func processLinePairs(filename string, allowLeftover bool, parser dualParser) (err error) {
	// load the file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// read a file line by line
	scanner := bufio.NewScanner(file)
	//adjust the capacity to your need (max characters in line)
	buf := make([]byte, bufferCapcityInBytes*1024)
	scanner.Buffer(buf, bufferCapcityInBytes*1024)

	readFirstLine := false
	var firstLine string
	for scanner.Scan() {
		// we have to read the first line first
		if !readFirstLine {
			firstLine = scanner.Text()
			readFirstLine = true

			// we read the first one already, so read the second one
		} else {
			err := parser(firstLine, scanner.Text())
			if err != nil {
				return err
			}

			firstLine = ""
			readFirstLine = false
		}
	}

	if readFirstLine && !allowLeftover {
		return errors.New("File did not contain an even number of lines")
	}

	// if something broke, throw an error
	return scanner.Err()
}

type iterCallback func(string) error

func iterateFiles(dir string, extension string, callback iterCallback) (err error) {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == extension {
			return callback(path)
		}
		return nil
	})
}
