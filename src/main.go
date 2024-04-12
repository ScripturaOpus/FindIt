package main

import (
	"findit/argparse"
	"findit/dirsearch"
	"findit/verbose"
	"fmt"
)

func main() {
	// Parse args into argparse.Config
	argparse.ParseArgs()

	verbose.Verbose("options", fmt.Sprintf("%v\n", argparse.Config))

	found := dirsearch.StartSearch()
	fmt.Println("----------------------------")
	fmt.Printf("Items found: %d\n", found)
}
