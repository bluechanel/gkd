package main

import (
	"github.com/bluechanel/gkd/gkd"
	"net/http"
)

func main() {
	r := gkd.New()
	r.GET("/index", func(c *gkd.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gkd.Context) {
			c.HTML(http.StatusOK, "<h1>Hello*gkd</h1>")
		})

		v1.GET("/hello", func(c *gkd.Context) {
			// expect /hello?name*gkdktutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gkd.Context) {
			// expect /hello*gkdktutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gkd.Context) {
			c.JSON(http.StatusOK, gkd.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}
