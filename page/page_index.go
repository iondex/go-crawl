package page

import (
	"github.com/go-redis/redis"
	"github.com/iondex/scraper-go/config"
	"github.com/iondex/scraper-go/requester"
	log "github.com/sirupsen/logrus"
)

var (
	ilogger = log.WithField("module", "pages")
)

// Index object is used to index pages stored in redis.
type Index struct {
	redisKey string
	client   *redis.Client
}

// NewIndex make a new url index and connect to redis server
func NewIndex() *Index {
	return &Index{
		redisKey: config.RedisPagesKey,
		client: redis.NewClient(&redis.Options{
			Addr:     config.RedisAddr,
			Password: config.RedisPassword,
		}),
	}
}

// Add is a temporary function. Will be removed later.
func (u *Index) Add(p *requester.Page) {
	b, err := u.client.HSet(config.RedisPagesKey, p.Url, p.Content).Result()
	if err != nil {
		ilogger.Errorf("AddPage failed - %s\n", err.Error())
		return
	}

	if b {
		ilogger.Infof("Page added: %s", p.Url)
	} else {
		ilogger.Warnf("Page already exists: %s", p.Url)
	}
}

// Len return length of the page index
func (u *Index) Len() int64 {
	n, err := u.client.HLen(u.redisKey).Result()
	if err != nil {
		ilogger.Errorf("Get page index length failed: %s", err.Error())
	}
	return n
}

// GetAllKeys returns all keys(urls) in the page index
func (u *Index) GetAllKeys() []string {
	a, err := u.client.HKeys(u.redisKey).Result()
	if err != nil {
		ilogger.Errorf("Get page index keys failed: %s", err.Error())
	}
	return a
}
