/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package main

import (
	"log"
	"os"

	"github.com/shinexia/elasticdump/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := cmd.NewElasticDumpCommand(os.Stdin, os.Stdout, os.Stderr)
	err := app.Execute()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}
