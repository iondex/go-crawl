package main

import (
	"flag"
	"strings"
	"time"

	"github.com/iondex/scraper-go/config"

	"github.com/iondex/scraper-go/page"
	"github.com/iondex/scraper-go/requester"
	log "github.com/sirupsen/logrus"
)

func parseFlags() {
	redisAddr := flag.String("redisAddr", config.RedisAddr, "Connection address of redis server.")
	redisPw := flag.String("redisPw", config.RedisPassword, "Redis server AUTH password.")
	maxConn := flag.Int("maxConn", config.MaxConcurrent, "Max concurrent goroutines for request.")
	reqInterval := flag.Int("sleep", config.RequesterInterval, "Interval between requests in milliseconds.")
	flag.Parse()

	config.RedisAddr = *redisAddr
	config.MaxConcurrent = *maxConn
	config.RedisPassword = *redisPw
	config.RequesterInterval = *reqInterval
}

func main() {
	parseFlags()
	logger := log.WithField("module", "main")

	r := requester.NewRequester(config.MaxConcurrent)
	p := page.NewParser()
	r.LinksIn(p.LinksOut())
	p.PagesIn(r.PagesOut())

	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	p.LinkFilter = func(url string) bool {
		return strings.HasPrefix(url, "http") && strings.Contains(url, "book.douban.com") && strings.Contains(url, "subject")
	}

	logger.Info("Crawler started.")
	r.AddTask("http://book.douban.com")
	time.Sleep(time.Hour)
	logger.Info("(Should) Shutdown.")
}
