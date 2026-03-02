package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
  g errgroup.Group
)

// routerServer 开启一个html服务器
func routerServer() http.Handler {
    e := gin.New()
    e.Use(gin.Recovery())
    e.LoadHTMLFiles("simple/index.html")
    e.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })
    return e
}

func cors1() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin","*")
        c.Next()
    }
}

func cors2(allowOrigin []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("vary", "origin")
        for _, o := range allowOrigin{
            if o == c.Request.Header.Get("origin"){
                c.Header("Access-Control-Allow-Origin", o)
                break
            }
        }

        c.Next()
    }
}

// routerApi 提供 api 服务
func routerApi() http.Handler {
    e := gin.New()
    e.Use(gin.Recovery())
    e.Use(cors2([]string{"http://127.0.0.1:4000", "http://localhost:4000"}))
    e.GET("/data", func(c *gin.Context) {
        c.String(200, "router api响应")
    })

    return e
}


func main() {
    server01 := &http.Server{
        Addr:         ":4000",
        Handler:      routerServer(),
    }

    server02 := &http.Server{
        Addr:         ":4001",
        Handler:      routerApi(),
    }

    g.Go(func() error {
        return server01.ListenAndServe()
    })

    g.Go(func() error {
        return server02.ListenAndServe()
    })

    if err := g.Wait(); err != nil {
        log.Fatal(err)
    }
}