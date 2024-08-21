package main

import (
	"log"

	"github.com/spf13/cobra/doc"

	"github.com/harvey-earth/elilogs/cmd"
)

func main() {
	title := &doc.GenManHeader{Title: "ELILOGS", Section: "1", Source: "harvey-earth", Manual: "elilogs Man Page"}
	err := doc.GenManTree(cmd.Root(), title, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
