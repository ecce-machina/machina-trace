package main

import (
	"fmt"
	"os"

	"github.com/ecce-machina/machina-trace/internal/diff"
	"github.com/ecce-machina/machina-trace/internal/render"
	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "diff":
		beforePath, afterPath, raw, err := parseDiffArgs(os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			usage()
			os.Exit(2)
		}
		runDiff(beforePath, afterPath, raw)

	default:
		usage()
		os.Exit(2)
	}
}

func parseDiffArgs(args []string) (beforePath, afterPath string, raw bool, err error) {
	var paths []string

	for _, arg := range args {
		switch {
		case arg == "--raw":
			raw = true

		case len(arg) > 0 && arg[0] == '-':
			return "", "", false, fmt.Errorf("unknown option: %s", arg)

		default:
			paths = append(paths, arg)
		}
	}

	if len(paths) != 2 {
		return "", "", false,
			fmt.Errorf("diff requires before.json and after.json")
	}

	return paths[0], paths[1], raw, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  machina-trace diff before.json after.json")
	fmt.Fprintln(os.Stderr, "  machina-trace diff --raw before.json after.json")
}

func runDiff(beforePath, afterPath string, raw bool) {
	before, err := snapshot.ReadFile(beforePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", beforePath, err)
		os.Exit(1)
	}

	after, err := snapshot.ReadFile(afterPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", afterPath, err)
		os.Exit(1)
	}

	deltas := diff.DiffSnapshots(before, after)

	render.WriteDiskFeaturesText(os.Stdout, deltas)
	render.WriteNetworkFeaturesText(os.Stdout, deltas)
	render.WriteLustreClientIOFeaturesText(os.Stdout, deltas)
	render.WriteLustreMetadataFeaturesText(os.Stdout, deltas)
	render.WriteLustreOSTFeaturesText(os.Stdout, deltas)

	if raw {
		render.WriteDiffText(os.Stdout, deltas)
	}

}
