package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageParser struct {
	input      chan string
	links      chan string
	linkFilter func(link string) bool
}

func NewPageParser(bufferLen int) *PageParser {
	r := &PageParser{
		links: make(chan string, bufferLen),
		input: make(chan string, bufferLen),
		linkFilter: func(url string) bool {
			return strings.HasPrefix(url, "http")
		},
	}
	r.Start()
	return r
}

// FlowIn starts a goroutine to read from chan in.
func (p *PageParser) FlowIn(in chan string) {
	go func() {
		for b := range in {
			p.input <- b
		}
	}()
}

func (p *PageParser) extractLinks(doc *goquery.Document) {
	sel := doc.Find("a")
	sel.Each(func(i int, sel *goquery.Selection) {

	})
}

func (p *PageParser) parse(body string) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return err
	}
	p.extractLinks(doc)
	return nil
}

func (p *PageParser) Start() {
	go func() {
		for body := range p.input {
			err := p.parse(body)
			if err != nil {
				log.Println("[Parser] Failed:", err)
			}
		}
	}()
}
