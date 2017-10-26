package server

import (
	"flag"
	"github.com/penguinn/penguin/component/config"
	"github.com/penguinn/penguin/component/db"
	"github.com/penguinn/penguin/component/log"
	"github.com/penguinn/penguin/component/router"
	"github.com/penguinn/penguin/component/mongo"
	"github.com/penguinn/penguin/component/redis"
)

func init() {
	flag.Parse()

	Use(config.ConfigComponent{})
	Use(log.LogComponent{})
	Use(router.RouterComponent{})
	Use(db.DBComponent{})
	Use(mongo.MongoComponent{})
	Use(redis.RedisComponent{})
}
