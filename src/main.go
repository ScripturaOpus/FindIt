package main

import (
	"findit/argparse"
	"findit/dirsearch"
	"fmt"
)

func main() {
	// Parse args into argparse.Config
	argparse.ParseArgs()

	fmt.Printf("%v\n", argparse.Config)

	found := dirsearch.StartSearch()
	fmt.Println("----------------------------")
	fmt.Printf("Items found: %d\n", found)
}
