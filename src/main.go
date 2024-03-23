/*
This is a learning project for understanding HTTP Caching Headers

# Caching Headers:
1. expires
2. pragma
3. content-control

# Validators:
1. etag
2. if-none-match
3. last-modified
4. if-modified-since
*/
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const PortEnv = "PORT"

type json map[string]interface{}
type jsonRouterHandler func(*gin.Context) (int, json)

func main() {
	fmt.Println("Starting server at port: ", os.Getenv(PortEnv))

	r := gin.Default()
	registerIndex(r)
	registerJsonRoutes(r)
	r.Run()
}

func registerIndex(r *gin.Engine) {
	r.LoadHTMLFiles("src/index.html")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})
}

func registerJsonRoutes(r *gin.Engine) {
	r.GET("/ping", jsonRouteWrapper(getPing))
	r.GET("/5_sec_expires", jsonRouteWrapper(_5SecExpires))
}

func jsonRouteWrapper(handler jsonRouterHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status, data := handler(ctx)
		ctx.JSON(status, data)
	}
}

func getPing(_ *gin.Context) (int, json) {
	return http.StatusOK, json{
		"message": "pong",
	}
}

// _5SecExpires: Creates an JSON with 5 seconds of expiry using the `Expires`
// response header. This header was available since before HTTP-1.1
func _5SecExpires(ctx *gin.Context) (int, json) {
	expires := getTime(5)
	ctx.Header("Expires", expires)

	return http.StatusOK, json{
		"message": "Using the `Expires` header setting the expiry of this JSON data blob to 5 seconds",
		"headers": json{
			"expires": expires,
		},
	}
}

func getTime(addSec int) string {
	now := time.Now().Add(time.Second * time.Duration(addSec)).UTC()
	return now.Format(time.RFC1123)
}
