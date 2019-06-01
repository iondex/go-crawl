package main

import (
	"time"

	"github.com/iondex/go-crawl/page"
	"github.com/iondex/go-crawl/requester"
)

func main() {
	r := requester.NewRequester(16)
	p := page.NewPageParser()
	r.In(p.Out())
	ign := r.Out()
	go func() {
		for {
			<-ign
		}
	}()

	r.AddTask("http://book.douban.com")
	time.Sleep(time.Hour)
}
