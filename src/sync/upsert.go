package sync

import (
	"fmt"
)

// upsertSegments updates or inserts new segements
func (proc *Process) upsertSegments() (err error) {
	fmt.Println("Upserting harvest segments ...")

	err = iterateFiles(proc.harvestFolder, ".json", func(path string) error {
		fmt.Printf("=> %s\n", path)
		return proc.upsertSegment(path)
	})

	if err == nil {
		fmt.Println("OK")
	} else {
		fmt.Println("ERROR")
	}

	return
}

// upsertSegment upserts a single segment
func (proc *Process) upsertSegment(segment string) (err error) {
	// compute the hash
	fmt.Print("  computing hash ... ")
	hash, err := proc.hashSegment(segment)

	if err != nil {
		fmt.Println("ERROR")
		return err
	}
	fmt.Printf("%s\n", hash)

	// check the index
	fmt.Print("  checking segment index ... ")

	obj, created, err := proc.checkSegmentIndex(segment)
	if err != nil {
		fmt.Println("ERROR")
		return err
	}

	if created {
		fmt.Println("NOT FOUND")
	} else {
		fmt.Println("FOUND")
	}

	// if the hash matches, we don't need to update
	if obj.Fields["hash"].(string) != hash {
		fmt.Println("  Hash in database differs")

		fmt.Print("  Clearing harvests belonging to segment ... ")
		err = proc.clearSegmentHarvests(segment)
		if err != nil {
			fmt.Println("ERROR")
			return err
		}
		fmt.Println("OK")

		// we need to clear out the old segments from the db, and put the new ones in
		fmt.Print("  Loading harvests from segment into index ... ")
		err = proc.insertSegmentHarvests(segment)
		if err != nil {
			fmt.Println("ERROR")
			return err
		}
		fmt.Println("OK")
	} else {
		fmt.Println("  Hash in database matches")
	}

	fmt.Print("  Storing segment state ... ")
	obj.Fields["touched"] = true
	obj.Fields["hash"] = hash

	// save it
	err = obj.Save()
	if err == nil {
		fmt.Println("OK")
	} else {
		fmt.Println("ERROR")
	}

	return err
}
