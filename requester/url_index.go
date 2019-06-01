package requester

import (
	"log"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-redis/redis"
	"github.com/iondex/go-crawl/config"
)

// UrlIndex object is used to prevent accessing previously accessed url.
type UrlIndex struct {
	data     mapset.Set
	redisKey string
	client   *redis.Client
}

// NewUrlIndex make a new url index and connect to redis server
func NewUrlIndex(key string) *UrlIndex {
	return &UrlIndex{
		data:     mapset.NewSet(),
		redisKey: key,
		client: redis.NewClient(&redis.Options{
			Addr:     config.RedisAddr,
			Password: config.RedisPassword,
		}),
	}
}

// Add adds a new url in UrlIndex. Errors won't be returned, but will be logged.
func (u *UrlIndex) Add(url string) {
	n, err := u.client.SAdd(u.redisKey, url).Result()
	if err != nil {
		log.Printf("WARNING: Add action failed - %s\n", err.Error())
	}
	if n != 1 {
		log.Printf("WARNING: Add action seems failed (n=%d)\n", n)
	}
}

// Has detects if url is in the UrlIndex
func (u *UrlIndex) Has(url string) (bool, error) {
	b, err := u.client.SIsMember(u.redisKey, url).Result()
	if err != nil {
		return true, err // don't add anything if error occurs
	}
	return b, nil
}

// AddPage is a temporary function. Will be removed later.
func (u *UrlIndex) AddPage(p *Page) {
	m := make(map[string]interface{})
	m[p.url] = p.content
	r, err := u.client.HMSet("pages", m).Result()
	if err != nil {
		log.Printf("WARNING: AddPage failed - %s\n", err.Error())
	} else {
		log.Println(r)
	}
}
