package main

import (
	"flag"
	"fmt"

	"github.com/Tantatorn-dev/gindle/pkg/crawler"
)

func main() {
	cookie := flag.String("cookie", "", "Cookie from read.amazon.com")
	targetPath := flag.String("target", "./notes", "Target path for JSON file")

	flag.Parse()

	if *cookie == "" {
		panic("Cookie is required")
	}

	bc := crawler.GetBooks(*cookie)

	for _, b := range bc {
		fmt.Println(b.Title, b.Author, len(b.Notes))
		_, err := b.ToJson(*targetPath)
		if err != nil {
			panic(err)
		}
	}
}
