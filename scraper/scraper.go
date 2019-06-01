package scraper

import (
	"github.com/iondex/scraper-go/page"
	"github.com/iondex/scraper-go/requester"
)

type Scraper struct {
	requester requester.Requester
	index     page.Index
}
