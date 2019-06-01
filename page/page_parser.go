package page

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/iondex/go-crawl/config"
)

var (
	maxChannelLen = config.MaxChannelLen
)

type PageParser struct {
	input      chan string
	links      chan string
	linkFilter func(link string) bool
}

func NewPageParser() *PageParser {
	r := &PageParser{
		links: make(chan string, maxChannelLen),
		input: make(chan string, maxChannelLen),
		linkFilter: func(url string) bool {
			return strings.HasPrefix(url, "http")
		},
	}
	r.Start()
	return r
}

// In starts a goroutine to read from chan in.
func (p *PageParser) In(in chan string) {
	go func() {
		for b := range in {
			p.input <- b
		}
	}()
}

// Out starts a goroutine to output all parsed links.
func (p *PageParser) Out() chan string {
	o := make(chan string, maxChannelLen)
	go func() {
		for l := range p.links {
			o <- l
		}
	}()
	return o
}

func (p *PageParser) extractLinks(doc *goquery.Document) {
	sel := doc.Find("a")
	sel.Each(func(i int, s *goquery.Selection) {
		v, e := s.Attr("href")
		if e && p.linkFilter(v) {
			p.links <- v
		}
	})
}

func (p *PageParser) parse(body string) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return err
	}
	go p.extractLinks(doc)
	return nil
}

// Start starts PageParser goroutine
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
