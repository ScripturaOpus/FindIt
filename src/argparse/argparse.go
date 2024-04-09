package argparse

import (
	"fmt"

	"findit/version"

	"github.com/alexflint/go-arg"
)

type options struct {

	/* Positional */
	SearchString     string `arg:"positional" help:"The string FindIt will look for in files"`
	WorkingDirectory string `arg:"positional" default:"." help:"The directory for FindIt to look in"`

	/* Switches */
	FileType string `arg:"-t:--type" help:"Specify a file type (or extension)"`

	/* Debugging */
	Verbose bool `arg:"--verbose" default:"false"`
}

var Config options

func (options) Version() string {
	return fmt.Sprintf("v%s\nch:%s\nbt:%s\n", version.VERSION, version.COMMIT_HASH, version.BUILD_TIME)
}

func ParseArgs() {
	_ = arg.MustParse(&Config)
}
