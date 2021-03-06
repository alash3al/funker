package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Debug = true
	e.HideBanner = true

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 9}))
	e.Use(middleware.Recover())
	e.Use(authMiddleware)

	e.GET("/", routeHome)
	e.POST("/funk/add/:funkName", routeAddFunk)
	e.DELETE("/funk/delete/:funkName", routeDeleteFunk)
	e.Any("/funk/exec/:funkName", routeCallFunk)
	e.Any("/funk/exec", routeEvalFunk)

	log.Fatal(e.Start(*flagListenAddr))
}
