package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var g errgroup.Group

func routerServer() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.LoadHTMLFiles("preflight/index.html")
	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	return e
}


func corsPreflight1() gin.HandlerFunc {
    return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin","*")
		if c.Request.Method == http.MethodOptions{
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

        c.Next()
    }
}

func corsPreflight2(allowOrigin []string) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Vary", "Origin")
		c.Header("Vary", "Access-Control-Request-Method")
		c.Header("Vary", "Access-Control-Request-Headers")

		requestOrigin := c.GetHeader("Origin")
		isAllowed := false

		for _, o := range allowOrigin {
			if o == requestOrigin {
				c.Header("Access-Control-Allow-Origin", o)
				isAllowed = true
				break
			}
		}

		if isAllowed && c.Request.Method == http.MethodOptions {
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type") // 允许前端带 Content-Type 头
			c.Header("Access-Control-Max-Age", "86400")              // 缓存 24 小时

			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func routerApi() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	e.Use(corsPreflight1())
	// e.Use(corsPreflight2([]string{"http://127.0.0.1:4000", "http://localhost:4000"}))

	e.POST("/data", func(c *gin.Context) {
		c.String(200, "router api响应：成功处理了复杂请求！")
	})

	return e
}

func main() {
	server01 := &http.Server{Addr: ":4000", Handler: routerServer()}
	server02 := &http.Server{Addr: ":4001", Handler: routerApi()}

	g.Go(func() error { return server01.ListenAndServe() })
	g.Go(func() error { return server02.ListenAndServe() })

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}