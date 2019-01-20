// Package supportscolor provides information on the level of color support
// found in a given terminal, obeying OS-specific color flags.
// Testing currently has been done via. iTerm2 manual testing.
package supportscolor

// TODO: Unit Tests, Windows Testing + Different popular terminal testing.

import (
	"os"
	"strconv"

	// TODO: Remove this later.
	"github.com/mattn/go-isatty"
)

var (
	colorFlags   = []string{"colors", "color=true", "color=always"}
	noColorFlags = []string{"no-color", "no-colors", "color=false"}
	cistrings    = []string{"TRAVIS", "CIRCLECI", "APPVEYOR", "GITLAB_CI"}

	flagMap = make(map[string]bool)

	// ForcedColor determines if there are flags in place forcing color.
	// This should never really change during the lifetime of a given program.
	ForcedColor *bool
)

// init() in this case simply adds the flags to a given map to make lookup easier.
func init() {
	for _, elem := range os.Args[1:] {
		flagMap[elem] = true
	}
}

// ColorSupport represents the given colors the terminal supports.
type ColorSupport struct {
	Level    int
	HasBasic bool
	Has256   bool
	Has1m    bool
}

// checkForcedColor checks to see if the terminal color is set
// to be explicitly a given support level.
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

// supportsColor finds the given support level of the temrinal,
// and returns it as an int representing a given level of support.
// TODO: Turn that level into a semi-enum via. iota.
func supportsColor() int {
	if ForcedColor == nil {
		checkForcedColor()
	}

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

	// Check to see if our terminal is truecolour
	if res, suc := os.LookupEnv("COLORTERM"); suc && res == "truecolor" {
		return 3
	}

	return min()
}

// TranslateLevel translates the given integer support level into a struct
// holding the general levels of support available in the given terminal.
func translateLevel(c int) *ColorSupport {
	var support = &ColorSupport{Level: c}
	switch {
	case c >= 1:
		support.HasBasic = true
		fallthrough
	case c >= 2:
		support.Has256 = true
		fallthrough
	case c == 3:
		support.Has1m = true
		break
	}
	return support
}

// GetSupportLevel returns a struct containing whether or not certain
// color sizes are supported in a given terminal.
func GetSupportLevel() *ColorSupport {
	if ForcedColor == nil {
		checkForcedColor()
	}
	return translateLevel(supportsColor())
}
