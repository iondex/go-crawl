package page

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/iondex/scraper-go/config"
	"github.com/iondex/scraper-go/requester"
	log "github.com/sirupsen/logrus"
)

var (
	maxChannelLen = config.MaxChannelLen
	logger        = log.WithField("module", "page_parser")
)

// Parser page parser struct
type Parser struct {
	input      chan *requester.Page
	links      chan string
	LinkFilter func(link string) bool
}

func NewParser() *Parser {
	p := &Parser{
		links: make(chan string, maxChannelLen),
		input: make(chan *requester.Page, maxChannelLen),
		LinkFilter: func(url string) bool {
			return strings.HasPrefix(url, "http")
		},
	}
	p.Start()
	return p
}

// PagesIn starts a goroutine to read pages from chan in.
func (p *Parser) PagesIn(in chan *requester.Page) {
	go func() {
		for b := range in {
			p.input <- b
		}
	}()
}

// LinksOut starts a goroutine to output all parsed links.
func (p *Parser) LinksOut() chan string {
	out := make(chan string, maxChannelLen)
	go func() {
		for l := range p.links {
			out <- l
		}
	}()
	return out
}

func (p *Parser) extractLinks(doc *goquery.Document) int {
	sel := doc.Find("a")
	sel.Each(func(i int, s *goquery.Selection) {
		v, e := s.Attr("href")
		if e && p.LinkFilter(v) {
			p.links <- v
		}
	})
	return len(sel.Nodes)
}

func (p *Parser) parse(page *requester.Page) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page.Content))
	if err != nil {
		return err
	}
	n := p.extractLinks(doc)
	logger.Infof("Extracted %d links from %s", n, page.Url)
	return nil
}

// Start starts PageParser goroutine
func (p *Parser) Start() {
	go func() {
		for page := range p.input {
			err := p.parse(page)
			if err != nil {
				log.Printf("[Parser] Failed for [%s] - %s\n", page.Url, err.Error())
			}
		}
	}()
}
