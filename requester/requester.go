package requester

import (
	"compress/gzip"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/iondex/scraper-go/config"
	log "github.com/sirupsen/logrus"
)

var (
	maxChannelLen = config.MaxChannelLen
	logger        = log.WithField("module", "requester")
)

// Requester is a global request client. Not thread safe.
type Requester struct {
	client   *http.Client
	jobs     chan string
	buffer   chan *Page
	size     int
	urlIndex *UrlIndex
}

// Page is used to represent a page with url and body.
type Page struct {
	Url     string
	Content string
}

// NewRequester construct a new Requester and start it.
func NewRequester(size int) *Requester {
	r := &Requester{
		client:   &http.Client{},
		jobs:     make(chan string, size),
		buffer:   make(chan *Page, maxChannelLen),
		size:     size,
		urlIndex: NewUrlIndex(config.RedisUrlIndexKey),
	}
	r.Start()
	return r
}

// AddTask add a task to Requester queue;
// the call might be blocked if channel is full.
func (r *Requester) AddTask(url string) {
	r.jobs <- url
}

func isHTML(contentType string) bool {
	if contentType == "" {
		return false
	}

	types := strings.Split(contentType, ",")
	for _, v := range types {
		t, _, _ := mime.ParseMediaType(v)
		if t == "text/html" {
			return true
		}
	}
	return false
}

// the place to perform real requests
func (r *Requester) request(url string) (*Page, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range config.DefaultHeaders {
		req.Header.Set(k, v)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP code %d: %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !isHTML(contentType) {
		return nil, errors.New("Target is not html")
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
		defer reader.Close()
	case "deflate":
		reader, _ = zlib.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return &Page{Url: url, Content: string(body)}, nil
}

// Start starts all worker goroutines.
func (r *Requester) Start() {
	for n := 0; n < r.size; n++ {

		go func(i int, jobs chan string) {
			workerLogger := logger.WithField("worker", i)

			for url := range r.jobs {
				page, err := r.request(url)
				if err != nil {
					workerLogger.Warnf("Request Error: %s\n", err.Error())
					time.Sleep(50 * time.Millisecond)
				} else {
					workerLogger.Infof("Request Done: %s\n", url)
					r.buffer <- page
					r.urlIndex.Add(url)
					r.urlIndex.AddPage(page)
				}
			}
		}(n, r.jobs)
	}
}

// Out creates arrays of output channels.
// Notice that by default the returned chan is buffered to facilitate async tasking.
func (r *Requester) PagesOut() chan *Page {
	out := make(chan *Page, maxChannelLen)
	go func() {
		for b := range r.buffer {
			out <- b
		}
	}()
	return out
}

// LinksIn accept input url channel.
func (r *Requester) LinksIn(in chan string) {
	go func() {
		for i := range in {
			b, err := r.urlIndex.Has(i)
			if err != nil {
				log.WithField("module", "requester").Printf("UrlIndex.Has failed - %s\n", err.Error())
			}
			// Has will return true when error occured, to prevent recursive
			if !b {
				r.jobs <- i
			}
		}
	}()
}
