package sync

import (
	"fmt"

	"github.com/olivere/elastic"
)

// Process represent args to the syncronisation process
type Process struct {
	client        *elastic.Client
	harvestFolder string

	harvestIndex string
	harvestType  string

	segmentIndex string
	segmentType  string
	segmentField string
}

// NewProcess creates a new Process
func NewProcess(Client *elastic.Client, Index string, Folder string) *Process {
	return &Process{
		client: Client,

		// harvests
		harvestIndex:  Index,
		harvestType:   "_doc",
		harvestFolder: Folder,

		// segments
		segmentIndex: fmt.Sprintf("%s-segments", Index),
		segmentType:  "_doc",
		segmentField: "segment",
	}
}

// Run is the main sync entry point
func (proc *Process) Run() {
	// Create the index and mapping
	err := proc.createIndex()
	if err != nil {
		panic(err)
	}

	// Reset the segment index
	err = proc.resetSegmentIndex()
	if err != nil {
		panic(err)
	}

	// refresh all the indexes
	err = proc.refreshIndex()
	if err != nil {
		panic(err)
	}

	// upsert segments
	err = proc.upsertSegments()
	if err != nil {
		panic(err)
	}

	// refresh all the indexes
	err = proc.refreshIndex()
	if err != nil {
		panic(err)
	}

	// clear old segements
	err = proc.clearSegments()
	if err != nil {
		panic(err)
	}

	// flush all the indexes
	err = proc.flushIndex()
	if err != nil {
		panic(err)
	}

	// and be done
	return
}
