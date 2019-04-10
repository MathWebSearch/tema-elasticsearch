package sync

import (
	"github.com/MathWebSearch/tema-elasticsearch/src/db"
)

// createIndex creates an index to prepare for segmented updates
func (proc *Process) createIndex() (err error) {
	proc.printf("Creating Harvest Index %s ... ", proc.harvestIndex)
	created, err := db.CreateIndex(proc.client, proc.harvestIndex, proc.makeHarvestMapping())
	if err != nil {
		proc.printlnERROR("ERROR")
		return
	}
	if created {
		proc.printlnOK("OK")
	} else {
		proc.printlnOK("SKIP")
	}

	proc.printf("Creating Segment Index %s ... ", proc.segmentIndex)
	created, err = db.CreateIndex(proc.client, proc.segmentIndex, proc.makeSegmentMapping())
	if err != nil {
		proc.printlnERROR("ERROR")
		return
	}
	if created {
		proc.printlnOK("OK")
	} else {
		proc.printlnOK("SKIP")
	}

	return
}

type l map[string]interface{}

// makeHarvestMapping makes a mapping for the harvest index
func (proc *Process) makeHarvestMapping() interface{} {
	return l{
		"settings": l{
			"index": l{
				"number_of_replicas": 0,
				"number_of_shards":   128,
			},
		},
		"mappings": l{
			proc.harvestType: l{
				"dynamic": false,
				"properties": l{
					"metadata": l{
						"dynamic": true,
						"type":    "object",
					},
					proc.segmentField: l{
						"type": "keyword",
					},
					"mws_ids": l{
						"type":  "long",
						"store": false,
					},
					"text": l{
						"type":  "text",
						"store": false,
					},
					"mws_id": l{
						"enabled": false,
						"type":    "object",
					},
					"math": l{
						"enabled": false,
						"type":    "object",
					},
				},
			},
		},
	}
}

// makeSegmentMapping makes a mapping for the segment index
func (proc *Process) makeSegmentMapping() interface{} {
	return l{
		"settings": l{
			"index": l{
				"number_of_replicas": 0,
				"number_of_shards":   128,
			},
		},
		"mappings": l{
			proc.segmentType: l{
				"dynamic": false,
				"properties": l{
					proc.segmentField: l{
						"type": "keyword",
					},
					"hash": l{
						"type": "keyword",
					},
					"touched": l{
						"type": "boolean",
					},
				},
			},
		},
	}
}
