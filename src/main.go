package main

import (
	"findit/argparse"
	"fmt"
)

func main() {
	// Parse args into argparse.Config
	argparse.ParseArgs()

	fmt.Printf("%v\n", argparse.Config)
}
