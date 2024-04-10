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
)

// For testing multiple search strings in the same line:
// test1    test2    test3      test4

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

func isBinary(filename string) (bool, error) {

	out, err := exec.Command("file", "-b", "--mime-type", filename).Output()

	if err != nil {
		return false, err
	}

	s_out := string(out)
	b_stat := false

	if strings.HasPrefix(s_out, "application") {
		b_stat = true
	}

	if argparse.Config.Verbose { // Just so we're not interpolating a string when we don't need to
		s_out = s_out[:len(s_out)-1]
		verbose.VerboseAny("return", fmt.Sprintf("%v, %s, %s", b_stat, filename, s_out))
	}

	return b_stat, nil
}

func matchesExt(path string) bool {
	if argparse.Config.FileType == "*" {
		return true
	}

	ext := filepath.Ext(path)
	reqType := strings.TrimPrefix(argparse.Config.FileType, "*")
	reqType = strings.TrimSuffix(reqType, "*")

	if ext == reqType {
		return true
	}

	// Check for wildcard prefix
	if strings.HasPrefix(ext, reqType) {
		return true
	}

	// Check for wildcard suffix
	if strings.HasSuffix(ext, reqType) {
		return true
	}

	return false
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
