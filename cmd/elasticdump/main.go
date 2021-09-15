package main

import (
	"log"
	"os"

	"github.com/shinexia/elasticdump/cmd/elasticdump/app"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd := app.NewElasticDumpCommand(os.Stdin, os.Stdout, os.Stderr)
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}
