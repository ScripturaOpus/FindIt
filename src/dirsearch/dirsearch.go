package dirsearch

import (
	"bufio"
	"findit/argparse"
	"findit/verbose"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sarpdag/boyermoore"
)

var foundItems uint64 // You never know. Could somehow be used.

/* Colors */
var file_descriptor_color string = "\033[36m"
var search_pattern_color string = "\033[31m"
var reset_color string = "\033[0m"

func StartSearch() uint64 {
	switch argparse.Config.ColorMode {
	case "auto":
		if !ansiSupported() {
			file_descriptor_color = ""
			reset_color = ""
			search_pattern_color = ""
		} else {
			argparse.Config.ColorEnabled = true
		}

	case "always":
		// Do nothing; Color is already setup
		argparse.Config.ColorEnabled = true

	case "never":
		file_descriptor_color = ""
		reset_color = ""
		search_pattern_color = ""
		argparse.Config.ColorEnabled = false

	default:
		fmt.Printf("Unknown color setting '%s'\n", argparse.Config.ColorMode)
		os.Exit(2)
	}

	recurseSearch(argparse.Config.WorkingDirectory)
	return foundItems
}

func recurseSearch(path string) {
	verbose.Verbose("Opening directory", path)
	entries, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("Failed to open directory:", err)
		return
	}

	for _, file := range entries {
		filePath := filepath.Join(path, file.Name())

		if file.IsDir() {
			// Scan the directory directory
			recurseSearch(filePath)
			continue
		}

		_ = checkFile(filePath, argparse.Config.SearchString)
	}
}

func checkFile(path string, search string) error {

	if !matchesExt(path) {
		return nil
	}

	isBin, err := isBinary(path)

	if err != nil {
		fmt.Println("Error detecting file:", err)
	}

	if (argparse.Config.OnlyBinaryFiles && !isBin) ||
		(!argparse.Config.AllowBinaryFiles && isBin) {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}

	verbose.Verbose("opening file", file.Name())

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1024*1024), -1)

	lineNum := 1
	searchLen := len(search)

	for scanner.Scan() {
		line := scanner.Text()
		lineLen := len(line)

		if line == "" || lineLen < searchLen { // The line cannot have what we're looking for if it's not big enough to contain it
			continue
		}

		var index int //boyermoore.Index(line, search)
		lastIndex := 0

		for {
			//verbose.Verbose("line", line[startIndex:])
			index = boyermoore.Index(line[lastIndex:], search)

			if index == -1 {
				break
			}

			index += lastIndex + 1

			if lastIndex > lineLen {
				break // Reached end of line, move to next line
			}

			isolated := isolate(line, search, argparse.Config.ContextSize, index-1)
			fmt.Printf("%s%s:%d:%d%s: %s\n", file_descriptor_color, path, lineNum, index-1, reset_color, isolated)
			foundItems++

			lastIndex = index
		}

		// Increment line number after processing the line
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	file.Close()
	return nil
}

// Get contextLength number of characters on either side of the substring
// If enabled, the substring will be red
func isolate(inputString string, substring string, contextLength int, startIndex int) string {
	index := strings.Index(substring, inputString[startIndex:])

	// Adjust startIndex and endIndex based on the actual found index (relative to startIndex)
	adjustedStartIndex := max(0, startIndex+index-contextLength)
	endIndex := min(len(inputString), startIndex+index+len(substring)+contextLength)

	// Replace substring with search pattern color codes (if enabled)
	output := inputString[adjustedStartIndex:endIndex]

	if argparse.Config.ColorEnabled {
		output = inputString[adjustedStartIndex:startIndex]
		output += search_pattern_color + substring + reset_color
		output += inputString[startIndex+len(substring) : endIndex]
	}

	return output
}
