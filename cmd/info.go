package cmd

import (
	"os"
	"runtime"

	"github.com/valord577/clix"
)

// @author valor.

var infoCmd = &clix.Command{
	Name:    "info",
	Summary: "Print information.",

	Run: printInfo,
}

var (
	version  = "dev"
	datetime = "-"
)

func printInfo(_ *clix.Command, _ []string) error {
	info := "webdav " + version + " " + datetime + " " + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH + "\n"
	_, err := os.Stdout.Write([]byte(info))
	return err
}
