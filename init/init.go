package init

import (
	"flag"
	"github.com/penguinn/penguin/component/config"
	"github.com/penguinn/penguin/component/db"
	"github.com/penguinn/penguin/component/log"
	"github.com/penguinn/penguin/component/router"
	"github.com/penguinn/penguin/component/server"
)

func init() {
	flag.Parse()

	server.Use(config.ConfigComponent{})
	server.Use(log.LogComponent{})
	server.Use(router.RouterComponent{})
	server.Use(db.DBComponent{})
}
