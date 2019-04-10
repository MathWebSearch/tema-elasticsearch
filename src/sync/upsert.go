package sync

import (
	"fmt"
)

// upsertSegments updates or inserts new segements
func (proc *Process) upsertSegments() (err error) {
	proc.println("Upserting harvest segments ...")

	err = iterateFiles(proc.harvestFolder, ".json", func(path string) error {
		proc.printf("=> %s\n", path)
		return proc.upsertSegment(path)
	})

	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return
}

// upsertSegment upserts a single segment
func (proc *Process) upsertSegment(segment string) (err error) {
	// compute the hash
	proc.print("  computing hash ... ")
	hash, err := proc.hashSegment(segment)

	if err != nil {
		proc.printlnERROR("ERROR")
		return err
	}
	proc.printf("%s\n", hash)

	// check the index
	proc.print("  checking segment index ... ")

	obj, created, err := proc.checkSegmentIndex(segment)
	if err != nil {
		proc.printlnERROR("ERROR")
		return err
	}

	if created {
		proc.printlnOK("NOT FOUND")
	} else {
		proc.printlnOK("FOUND")
	}

	proc.print("  Comparing segment hash ... ")

	// if the hash matches, we don't need to update
	if obj.Fields["hash"].(string) != hash {
		proc.printlnOK("DIFFERS")

		proc.print("  Clearing harvests belonging to segment ... ")
		err = proc.clearSegmentHarvests(segment)
		if err != nil {
			proc.printlnERROR("ERROR")
			return err
		}
		proc.printlnOK("OK")

		// we need to clear out the old segments from the db, and put the new ones in
		fmt.Print("  Loading harvests from segment into index ... ")
		err = proc.insertSegmentHarvests(segment)
		if err != nil {
			proc.printlnERROR("ERROR")
			return err
		}
		proc.printlnOK("OK")
	} else {
		proc.printlnOK("MATCHES")
	}

	proc.print("  Storing segment state ... ")
	obj.Fields["touched"] = true
	obj.Fields["hash"] = hash

	// save it
	err = obj.Save()
	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return err
}
