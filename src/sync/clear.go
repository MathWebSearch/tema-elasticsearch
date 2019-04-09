package sync

import (
	"fmt"

	"github.com/MathWebSearch/tema-elasticsearch/src/db"
	"github.com/olivere/elastic"
)

// clearSegments clears untouched (old) segments from the index
func (proc *Process) clearSegments() (err error) {
	fmt.Println("Removing old segments ...")

	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery("touched", false))

	old, err := db.FetchObjects(proc.client, proc.segmentIndex, proc.segmentType, q)
	if err != nil {
		for so := range old {
			e := proc.clearSegment(so)
			if e != nil {
				err = e
			}
		}
	}

	if err == nil {
		fmt.Println("OK")
	} else {
		fmt.Println("ERROR")
	}

	return
}

// clearSegment removes a single segment
func (proc *Process) clearSegment(so *db.ECObject) (err error) {
	segment := so.Fields[proc.segmentField].(string)
	fmt.Printf("=> %s\n", segment)

	// clear the harvests
	fmt.Print("  Clearing harvests belonging to segment ... ")
	err = proc.clearSegmentHarvests(segment)
	if err != nil {
		fmt.Println("ERROR")
		return
	}
	fmt.Println("OK")

	// and remove segment itself
	fmt.Print("  Removing segment ...")
	err = so.Delete()
	if err == nil {
		fmt.Println("OK")
	} else {
		fmt.Println("ERROR")
	}

	return

}

// clearSegmentHarvests clears segments belonging to a harvest
func (proc *Process) clearSegmentHarvests(segment string) error {
	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery(proc.segmentField, segment))
	return db.DeleteBulk(proc.client, proc.harvestIndex, proc.harvestType, q)
}
