package argparse

import (
	"fmt"
	"os"

	"findit/version"

	"github.com/alexflint/go-arg"
)

type options struct {

	/* Positional */
	SearchString     string `arg:"positional" help:"The string FindIt will look for in files within a directory"`
	WorkingDirectory string `arg:"positional" default:"." help:"The directory for FindIt to look in"`

	/* Switches */
	FileType         string `arg:"-t:--type" default:"*" help:"Specify a file type (or extension)"`
	ContextSize      int    `arg:"-c,--context-size" default:"20" help:"How many characters should be appended to the found strings context"`
	ColorMode        string `arg:"--color" default:"auto" help:"Specify if output should be colored (never, always, auto)"`
	AllowBinaryFiles bool   `arg:"-b,--allow-binary" default:"false" help:"Scan files detected as binary for the given search string (Not recommended)"`
	OnlyBinaryFiles  bool   `arg:"--only-binary" default:"false" help:"Only scan files detected as binary for the given search string (Not recommended)"`

	/* Debugging */
	Verbose bool `arg:"--verbose" default:"false"`
}

var Config options

/* Other */
var ColorEnabled bool

func (options) Version() string {
	if Config.SearchString == "" {
		return fmt.Sprintf("v%s", version.VERSION)
	}

	return fmt.Sprintf("v%s\nch:%s\nbt:%s", version.VERSION, version.COMMIT_HASH, version.BUILD_TIME)
}

func ParseArgs() {
	_ = arg.MustParse(&Config)

	if Config.SearchString == "" {
		// It's dumb that I even have to do this for help i