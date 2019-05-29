package main

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/go-redis/redis"
)

type UrlSet struct {
	data   mapset.Set
	key    string
	client *redis.Client
}

func NewUrlSet(key string) *UrlSet {
	return &UrlSet{
		data: mapset.NewSet(),
		key:  key,
	}
}

func (u *UrlSet) Sync() {
	if u.client == nil {
		u.client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "123456",
		})
	}
}
