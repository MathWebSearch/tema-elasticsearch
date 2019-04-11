package sync

import (
	"github.com/MathWebSearch/elasticutils"
)

func (proc *Process) refreshIndex() error {
	proc.print("Refreshing elasticsearch ... ")
	err := elasticutils.RefreshIndex(proc.client, proc.segmentIndex, proc.harvestIndex)

	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return err
}

func (proc *Process) flushIndex() error {
	proc.print("Flushing elasticsearch ... ")
	err := elasticutils.FlushIndex(proc.client, proc.segmentIndex, proc.harvestIndex)

	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return err
}
