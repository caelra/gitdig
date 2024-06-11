package main

import (
	"log"

	"github.com/caelra/gitdig/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
