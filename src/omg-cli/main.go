package main

import (
	"log"
	"os"

	"omg-cli/omg/commands"

	"github.com/alecthomas/kingpin"
)

func main() {
	logger := log.New(os.Stderr, "[OMG] ", 0)

	app := kingpin.New("omg-cli", "OMG! Ops Manager (on) Google")
	commands.Configure(logger, app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
