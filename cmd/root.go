package cmd

import (
	"github.com/valord577/clix"
)

// @author valor.

var rootCmd = &clix.Command{
	Name:    "webdav",
	Summary: "A Lightweight WebDAV Server.",
}

func init() {
	rootCmd.AddCmd(infoCmd, servCmd)
}
