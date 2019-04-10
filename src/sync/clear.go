package sync

import (
	"github.com/MathWebSearch/tema-elasticsearch/src/db"
	"github.com/olivere/elastic"
)

// clearSegments clears untouched (old) segments from the index
func (proc *Process) clearSegments() (err error) {
	proc.println("Removing old segments ...")

	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery("touched", false))

	old, err := db.FetchObjects(proc.client, proc.segmentIndex, proc.segmentType, q)
	if err == nil {
		for so := range old {
			e := proc.clearSegment(so)
			if e != nil {
				err = e
			}
		}
	}

	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return
}

// clearSegment removes a single segment
func (proc *Process) clearSegment(so *db.ECObject) (err error) {
	segment := so.Fields[proc.segmentField].(string)
	proc.printf("=> %s\n", segment)

	// clear the harvests
	proc.print("  Clearing harvests belonging to segment ... ")
	err = proc.clearSegmentHarvests(segment)
	if err != nil {
		proc.printlnERROR("ERROR")
		return
	}
	proc.printlnOK("OK")

	// and remove segment itself
	proc.print("  Removing segment ... ")
	err = so.Delete()
	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return

}

// clearSegmentHarvests clears segments belonging to a harvest
func (proc *Process) clearSegmentHarvests(segment string) error {
	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery(proc.segmentField, segment))
	return db.DeleteBulk(proc.client, proc.harvestIndex, proc.harvestType, q)
}
