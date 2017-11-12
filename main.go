package main

import (
	"flag"
	"github.com/penguinn/album/controllers"
	"github.com/penguinn/penguin/component/config"
	"github.com/penguinn/penguin/component/db"
	"github.com/penguinn/penguin/component/log"
	"github.com/penguinn/penguin/component/router"
	"github.com/penguinn/penguin/component/server"
	"github.com/penguinn/penguin/component/session"
)

func main() {
	flag.Parse()

	server.Use(config.ConfigComponent{})
	server.Use(log.LogComponent{})
	server.Use(router.RouterComponent{})
	server.Use(db.DBComponent{})

	router.RegisterController(controllers.NewBaseController(), session.Middleware)
	router.RegisterControllerGroup(controllers.NewCheckController(), "api", session.Middleware)
	router.RegisterControllerGroup(controllers.NewUserController(), "api", session.Middleware)
	router.RegisterControllerGroup(controllers.NewUploadController(), "api", session.Middleware)
	server.Serve()
}
