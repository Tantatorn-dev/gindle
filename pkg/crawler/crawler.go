package crawler

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

var TargetURL = "https://read.amazon.com/notebook"

const Cookie = `COOKIE`

type Book struct {
	ID     string
	Title  string
	Author string
	Notes  []string
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

	noteCollector := c.Clone()

	c.OnHTML("#kp-notebook-library", func(e *colly.HTMLElement) {
		attrs := e.ChildAttrs("div", "id")

		for i := range bs {
			bs[i].ID = attrs[i]
			bs[i].Notes = getNotes(noteCollector, attrs[i])
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", Cookie)
	})

	c.Visit(TargetURL)

	return bs
}

func getNotes(c *colly.Collector, id string) []string {
	var ns []string

	c.OnHTML(`div span[id=highlight]`, func(e *colly.HTMLElement) {
		ns = append(ns, e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", Cookie)
	})

	c.Visit(fmt.Sprintf("%s/?asin=%s&contentLimitState=&", TargetURL, id))

	return ns
}
