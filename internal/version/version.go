package version

import (
	"fmt"
	"os"

	"github.com/charlesnobrega/STARDF-Anime/internal/tracking"
)

var (
	Version   = "1.6.3"
	BuildTime = "unknown"
	Commit    = "unknown"
)

func HasVersionArg() bool {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		return arg == "--version" || arg == "-version" || arg == "-v" || arg == "--v"
	}
	return false
}

func ShowVersion() {
	fmt.Printf("StarDF-Anime v%s", Version)
	if tracking.IsCgoEnabled {
		fmt.Println(" (with SQLite tracking)")
	} else {
		fmt.Println(" (without SQLite tracking)")
	}

	if BuildTime != "" && BuildTime != "unknown" {
		fmt.Printf("Build time: %s\n", BuildTime)
	}
	if Commit != "" && Commit != "unknown" {
		fmt.Printf("Commit: %s\n", Commit)
	}
}
