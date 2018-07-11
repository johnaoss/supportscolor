package supportscolor

import (
	"os"
	"strconv"

	"github.com/mattn/go-isatty"
)

var (
	colorFlags   = []string{"colors", "color=true", "color=always"}
	noColorFlags = []string{"no-color", "no-colors", "color=false"}
	cistrings    = []string{"TRAVIS", "CIRCLECI", "APPVEYOR", "GITLAB_CI"}

	flagMap = make(map[string]bool)
	// ForcedColor determines if there are flags in place forcing colour.
	ForcedColor *bool
)

// Call before to make sure that it all runs smooth
func init() {
	for _, elem := range os.Args {
		flagMap[elem] = true
	}
}

// ColorSupport represents the... color the terminal supports
type ColorSupport struct {
	Level    int
	HasBasic bool
	Has256   bool
	Has1m    bool
}

func checkForcedColor() {
	var boolean bool
	if item, envcolor := os.LookupEnv("FORCE_COLOR"); envcolor {
		num, _ := strconv.Atoi(item)
		boolean = len(item) == 0 || num != 0
		ForcedColor = &boolean
		return
	}
	for _, elem := range colorFlags {
		if flagMap[elem] {
			boolean = true
			ForcedColor = &boolean
			return
		}
	}
	for _, elem := range noColorFlags {
		if flagMap[elem] {
			boolean = false
			ForcedColor = &boolean
			return
		}
	}
	ForcedColor = nil
}

func supportsColor() int {
	if ForcedColor != nil && *ForcedColor == false {
		return 0
	}

	if flagMap["color=16m"] || flagMap["color=full"] || flagMap["color=truecolor"] {
		return 3
	}

	if flagMap["color=256"] {
		return 2
	}

	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return 0
	}

	var min = func() int {
		if ForcedColor != nil {
			return 1
		}
		return 0
	}

	// Are we in a CI that supports colored text?
	if _, suc := os.LookupEnv("CI"); suc {
		for _, elem := range cistrings {
			if res, suc := os.LookupEnv(elem); suc && res != "" {
				return 1
			}
		}
		if res, suc := os.LookupEnv("CI_NAME"); suc && res == "codeship" {
			return 1
		}
		return min()
	}

	if res, suc := os.LookupEnv("COLORTERM"); suc && res == "truecolor" {
		return 3
	}

	return min()
}

/*
// GetSupportLevel determines if you can write color to the terminal.
// Currently only supports writing out to Files/Stdout/Stderr
func getSupportLevel() {
	const level = checkForcedColor()
	return translateLevel()
}

func translateLevel(...io.Writer) int {

}
*/
