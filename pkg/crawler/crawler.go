package crawler

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gocolly/colly/v2"
)

const TargetURL = "https://read.amazon.com/notebook"

type Book struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Author string   `json:"author"`
	Notes  []string `json:"notes"`
}

// ToJson is a method of for persisting a Book struct to JSON file
func (b Book) ToJson(targetPath string) (bool, error) {
	file, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return false, err
	}

	targetPath = path.Join(targetPath, fmt.Sprintf("%s.json", b.Title))

	err = os.WriteFile(targetPath, file, 0644)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetBooks is a function for getting all books from read.amazon.com
func GetBooks(cookie string) []Book {
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
			bs[i].Notes = getNotes(noteCollector, cookie, attrs[i])
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", cookie)
	})

	c.Visit(TargetURL)

	return bs
}

// getNotes is a function for getting all notes, highlight from a book
func getNotes(c *colly.Collector, cookie, id string) []string {
	var ns []string

	c.OnHTML(`div span[id=highlight]`, func(e *colly.HTMLElement) {
		ns = append(ns, e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", cookie)
	})

	c.Visit(fmt.Sprintf("%s/?asin=%s&contentLimitState=&", TargetURL, id))

	return ns
}
