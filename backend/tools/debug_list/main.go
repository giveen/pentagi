package main

import (
	"fmt"
	"log"

	"pentagi/cmd/installer/files"
)

func main() {
	f := files.NewFiles()
	list, err := f.List("observability")
	if err != nil {
		log.Fatalf("List error: %v", err)
	}
	fmt.Printf("found %d files:\n", len(list))
	for _, l := range list {
		fmt.Println(l)
	}
}
