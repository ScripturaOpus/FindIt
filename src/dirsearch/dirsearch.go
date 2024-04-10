package dirsearch

import (
	"bufio"
	"findit/argparse"
	"findit/verbose"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sarpdag/boyermoore"
)

// test    test    test    test

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
		fmt.Printf("Unknown color setting '%s'", argparse.Config.ColorMode)
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
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}

	verbose.Verbose("Opening file", file.Name())

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1024*1024), -1)

	lineNum := 0
	searchLen := len(search)

	for scanner.Scan() {
		line := scanner.Text()
		lineLen := len(line)

		index := boyermoore.Index(line, search)
		startIndex := 0

		for index != -1 {
			isolated := isolate(line, search, 20, index)

			fmt.Printf("%s%s:%d:%d%s %s\n", file_descriptor_color, path, lineNum, index, reset_color, isolated)
			foundItems++

			startIndex += index + searchLen
			if startIndex > lineLen {
				break // Reached end of line, move to next line
			}

			index = boyermoore.Index(line[startIndex:], search)
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

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func ansiSupported() bool {
	// Check using tput if available
	if _, err := exec.LookPath("tput"); err == nil {
		out, err := exec.Command("tput", "colors").Output()
		if err == nil {
			var numColors int // Declare numColors here

			_, err := fmt.Sscanln(string(out), &numColors)

			if err == nil && numColors >= 8 {
				return true
			}
		}
	}

	// Direct console query (ensure CSI is defined)
	const csi = "\033["
	fmt.Print(csi + "c")
	reader := bufio.NewReader(os.Stdin)
	ansiReport, err := reader.ReadString('c')

	if err != nil {
		return false
	}

	return strings.TrimSpace(ansiReport) != ""
}
