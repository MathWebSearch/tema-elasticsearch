package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/MathWebSearch/tema-elasticsearch/src/db"

	"github.com/olivere/elastic"
)

// checkSegmentIndex checks the segment index for a given segment
func (proc *Process) checkSegmentIndex(segment string) (obj *db.ECObject, created bool, err error) {
	// the query
	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery(proc.segmentField, segment))

	// the fields
	fields := make(map[string]interface{})
	fields[proc.segmentField] = segment
	fields["hash"] = ""
	fields["touched"] = true

	// and fetch or create it from the index
	return db.FetchOrCreateObject(proc.client, proc.segmentIndex, proc.segmentType, q, fields)
}

// hashSegment computes the hash of a segment
func (proc *Process) hashSegment(segment string) (hash string, err error) {
	// the hasher implementation
	hasher := sha256.New()

	// open the segment
	f, err := os.Open(segment)
	if err != nil {
		return
	}
	defer f.Close()

	// start hashing the file
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	// turn it into a string
	hash = hex.EncodeToString(hasher.Sum(nil))
	return
}
