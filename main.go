/*
Copyright © 2023 orangekame3 <miya.org.0309@gmai.com>
*/
package main

import (
	"github.com/orangekame3/tftarget/cmd"
)

var (
	version = "0.0.23"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
