package main

import (
	_ "github.com/penguinn/album/init"

	"github.com/penguinn/album/controllers"
	"github.com/penguinn/penguin/component/router"
	"github.com/penguinn/penguin/component/server"
	"github.com/penguinn/penguin/component/session"
)

func main() {

	router.Static("/static", "./templates")

	router.RegisterController(controllers.NewBaseController(), session.Middleware)
	router.RegisterControllerGroup(controllers.NewCheckController(), "api", session.Middleware)
	router.RegisterControllerGroup(controllers.NewUserController(), "api", session.Middleware)
	router.RegisterControllerGroup(controllers.NewUploadController(), "api", session.Middleware)
	server.Serve()

}
