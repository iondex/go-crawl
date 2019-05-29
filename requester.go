package main

import (
	"errors"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"
)

type Requester struct {
	client *http.Client
	jobs   chan string
	buffer chan string
	size   int
}

// NewRequester construct a new Requester and start it.
func NewRequester(size int) *Requester {
	r := &Requester{
		client: &http.Client{},
		jobs:   make(chan string, size),
		buffer: make(chan string, 1024),
		size:   size,
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

func (r *Requester) request(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !isHTML(contentType) {
		return errors.New("Target is not html")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	r.buffer <- string(body)
	return nil
}

// Start starts all worker goroutines.
func (r *Requester) Start() {
	for n := 0; n < r.size; n++ {
		go func(i int, jobs chan string) {
			for url := range r.jobs {
				err := r.request(url)
				if err != nil {
					log.Printf("[Worker %d] Error: %s\n", i, err)
					time.Sleep(50 * time.Millisecond)
				}
				log.Printf("[Worker %d] Done: %s\n", i, url)
			}
		}(n, r.jobs)
	}
}

// FlowOut creates arrays of output channels.
// Notice that by default the returned chan is buffered to facilitate async tasking.
func (r *Requester) FlowOut(size int) []chan string {
	out := make([]chan string, size)
	for i := 0; i < size; i++ {
		out[i] = make(chan string, 1024)
	}
	go func() {
		for b := range r.buffer {
			for _, c := range out {
				c <- b
			}
		}
	}()
	return out
}

func (r *Requester) FlowIn(in chan string) {
	go func() {
		for i := range in {
			r.jobs <- i
		}
	}()
}
