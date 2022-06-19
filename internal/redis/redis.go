package redis

import (
	"fmt"
	"github.com/go-redis/redis/v9"
	"os"
)

var rdb *redis.Client

const auctionStatePrefixKey = "auction_stats_"

func Init() error {
	host := os.Getenv("REDIS_URL")

	if host == "" {
		return fmt.Errorf("redis host env var not given")
	}

	opt, err := redis.ParseURL(host)
	if err != nil {
		return err
	}
	rdb = redis.NewClient(opt)

	return nil
}
