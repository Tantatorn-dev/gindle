package crawler

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

var TargetURL = "https://read.amazon.com/notebook"

var Cookie = `your_cookie`

type Book struct {
	Title  string
	Author string
	Note   []string
}

func GetBooks() []Book {
	var bs []Book

	c := colly.NewCollector()

	c.OnHTML("a h2", func(e *colly.HTMLElement) {
		b := Book{
			Title: e.Text,
		}

		bs = append(bs, b)
	})

	i := 0
	c.OnHTML("a p", func(e *colly.HTMLElement) {
		bs[i].Author = strings.Split(e.Text, ": ")[1]
		i++
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", Cookie)

		fmt.Println("Visiting", r.URL)
	})

	c.Visit(TargetURL)

	return bs
}
