package main

import (
	"github.com/bluechanel/gkd/gkd"
	"net/http"
)

func main() {
	r := gkd.New()
	r.GET("/", func(c *gkd.Context) {
		c.HTML(http.StatusOK, "<h1>hello gkd</h1>")
	})

	r.GET("/hello", func(c *gkd.Context) {
		c.String(http.StatusOK, "hello handler", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gkd.Context) {
		c.String(http.StatusOK,"hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("assets/*filepath", func(c *gkd.Context) {
		c.JSON(http.StatusOK, gkd.H{"filename":c.Param("filepath")})
	})
	r.Run(":9999")
}
