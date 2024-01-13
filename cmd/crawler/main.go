package main

import (
	"fmt"

	"github.com/Tantatorn-dev/gindle/pkg/crawler"
)

func main() {
	bc := crawler.GetBooks()

	for _, b := range bc {
		fmt.Println(b.ID, b.Title, b.Author)
	}
}
