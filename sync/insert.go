package sync

import (
	"encoding/json"

	"github.com/MathWebSearch/elasticutils"
)

// insertSegmentHarvests inserts a segment
func (proc *Process) insertSegmentHarvests(segment string) error {
	store := make(chan map[string]interface{})
	errChan := make(chan error)

	// start processing async
	go func() {
		e := processLinePairs(segment, true, func(_, contentLine string) (err error) {
			// unmarshal the content
			var content map[string]interface{}
			err = json.Unmarshal([]byte(contentLine), &content)
			if err != nil {
				return
			}

			// store the segment id
			content[proc.segmentField] = segment

			// and put it in the db
			store <- content

			return
		})

		// close both of the channel
		close(store)

		errChan <- e
		close(errChan)
	}()

	// run the insert and get the errors
	bulkError := elasticutils.CreateBulk(proc.client, proc.harvestIndex, proc.harvestType, store)
	parseError := <-errChan

	// return the parser error
	if parseError != nil {
		return parseError
	}

	// or the bulk error
	return bulkError
}
