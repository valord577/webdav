package main

import (
	"os"

	"github.com/valord577/webdav/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Stderr.Write([]byte("\n"))
		os.Exit(2)
	}
}
