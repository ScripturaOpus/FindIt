package verbose

import (
	"findit/argparse"
	"fmt"
	"runtime"
	"strings"
)

// Write message to the console when --verbose is used
func Verbose(prefix string, in string) string {
	if argparse.Config.Verbose {
		caller := CallerName(1)
		fmt.Printf("[%s] %s: `%s`\n", caller, prefix, in)
	}

	return in
}

func VerboseAny(prefix string, in any) {
	if argparse.Config.Verbose {
		caller := CallerName(1)
		fmt.Printf("[%s] %s: `%v`\n", caller, prefix, in)
	}
}

func CallerName(skip int) string {

	// Being someone who comes from C#, if feels weird being able to use
	// random empty quotes everywhere without having to use "string.Empty" for efficiency

	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}

	f := runtime.FuncForPC(pc)
	if f == nil {
		return ""
	}

	name := f.Name()

	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}
