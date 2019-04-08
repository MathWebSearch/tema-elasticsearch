package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/MathWebSearch/tema-elasticsearch/src"
)

func main() {
	// parse and validate arguments
	args := src.ParseArgs(os.Args)
	if !args.Validate() {
		die(nil)
	}

	// make a client and wait for es to come online
	client := src.MakeClientAndWait(args)

	// create the database or die
	err := client.SetupDB()
	if err != nil {
		die(err)
	}

	err = client.LoadHarvests()
	if err != nil {
		die(err)
	}

	fmt.Println("Done updating index")

}

func die(err error) {
	// KILL the parent process (inside docker)
	syscall.Kill(os.Getppid(), syscall.SIGTERM)

	if err != nil {
		panic(err)
	} else {
		panic("Something went wrong")
	}
}

// TODO: 1. Wait for tema-search to be up on the given port
// 2. Check if we have to run setup
// 3. Hash the directory; if it has changed clear out and fully re-index
