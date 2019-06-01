package requester

import (
	"github.com/go-redis/redis"
	"github.com/iondex/scraper-go/config"
	log "github.com/sirupsen/logrus"
)

// UrlIndex object is used to prevent accessing previously accessed url.
type UrlIndex struct {
	redisKey string
	client   *redis.Client
}

// NewUrlIndex make a new url index and connect to redis server
func NewUrlIndex(key string) *UrlIndex {
	return &UrlIndex{
		redisKey: key,
		client: redis.NewClient(&redis.Options{
			Addr:     config.RedisAddr,
			Password: config.RedisPassword,
		}),
	}
}

// Add adds a new url in UrlIndex. Errors won't be returned, but will be logged.
func (u *UrlIndex) Add(url string) {
	logger := log.WithFields(log.Fields{
		"module": "url_index",
	})
	n, err := u.client.SAdd(u.redisKey, url).Result()
	if err != nil {
		logger.Errorf("Add action failed: %s", err.Error())
	}
	if n == 0 {
		logger.Warnf("URL already exists: %s", url)
	} else {
		logger.Infof("URL added: %s", url)
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
	logger := log.WithField("module", "pages")
	b, err := u.client.HSet(config.RedisPagesKey, p.Url, p.Content).Result()
	if err != nil {
		logger.Errorf("AddPage failed - %s\n", err.Error())
		return
	}

	if b {
		logger.Infof("Page added: %s", p.Url)
	} else {
		logger.Warnf("Page already exists: %s", p.Url)
	}
}
