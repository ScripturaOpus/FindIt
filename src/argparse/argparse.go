package argparse

import (
	"fmt"
	"os"

	"findit/version"

	"github.com/alexflint/go-arg"
)

type options struct {

	/* Positional */
	SearchString     string `arg:"positional" help:"The string FindIt will look for in files"`
	WorkingDirectory string `arg:"positional" default:"." help:"The directory for FindIt to look in"`

	/* Switches */
	FileType    string `arg:"-t:--type" help:"Specify a file type (or extension)"`
	ContextSize int    `arg:"-c,--context-size" default:"20" help:"How many characters should be appended to the found strings context"`
	ColorMode   string `arg:"--color" default:"auto" help:"Specify if output should be colored (never, always, auto)"`

	/* Debugging */
	Verbose bool `arg:"--verbose" default:"false"`

	/* Other */
	ColorEnabled bool
}

var Config options

func (options) Version() string {
	return fmt.Sprintf("v%s\nch:%s\nbt:%s", version.VERSION, version.COMMIT_HASH, version.BUILD_TIME)
}

func ParseArgs() {
	_ = arg.MustParse(&Config)

	if Config.SearchString == "" {
		os.Args = []string{"", "--help"}
		arg.MustParse(&Config)
	}
}
