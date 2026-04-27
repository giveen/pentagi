package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	links := "./cmd/installer/files/links"
	info, err := os.Stat(links)
	if err != nil {
		log.Fatalf("stat links: %v", err)
	}
	fmt.Printf("links is dir: %v\n", info.IsDir())

	err = filepath.Walk(links, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(links, path)
		fmt.Printf("path: %s  dir:%v\n", rel, info.IsDir())
		return nil
	})
	if err != nil {
		log.Fatalf("walk error: %v", err)
	}
}
