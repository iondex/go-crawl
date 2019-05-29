package main

import (
	"log"
	"time"
)

func main() {
	r := NewRequester(2)
	cs := r.FlowOut(2)
	p := NewPageParser(16)
	go func(ch <-chan string) {
		for b := range ch {
			log.Println(len(b))
		}
	}(cs[0])
	p.FlowIn(cs[1])
	log.Println("Workers started")

	var urls = []string{
		"http://www.baidu.com/",
		"http://www.360.cn/",
		"http://www.qq.com/",
		"http://www.douban.com",
		"https://dl.360safe.com/inst.exe",
	}
	for _, u := range urls {
		r.AddTask(u)
	}

	time.Sleep(30 * time.Second)
}
