package sync

import (
	"github.com/MathWebSearch/elasticutils"
	"gopkg.in/olivere/elastic.v6"
)

// resetSegmentIndex resets the segment index to prepare for updates
func (proc *Process) resetSegmentIndex() (err error) {
	proc.print("Resetting Segment Index ... ")

	// reset the touched part to false
	err = elasticutils.UpdateAll(proc.client, proc.segmentIndex, proc.segmentType, elastic.NewScript("ctx._source.touched = false"))
	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return
}
