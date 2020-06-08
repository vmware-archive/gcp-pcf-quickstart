// +build tools

package main

import (
	"log"

	"omg-cli/templates"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(templates.Templates, vfsgen.Options{
		PackageName:  "templates",
		BuildTags:    "!dev",
		VariableName: "Templates",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
