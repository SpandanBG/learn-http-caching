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
	"math/rand"
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
	r.GET("/pragma", jsonRouteWrapper(_pragmaNoCache))
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
// response header. This header was available since before HTTP-1.1.
func _5SecExpires(ctx *gin.Context) (int, json) {
	expires := getTime(5)
	ctx.Header("Expires", expires)

	return http.StatusOK, json{
		"message": "Using the `Expires` header setting the expiry of this JSON data blob to 5 seconds.You can see the `random` value change if you request for it after 5 seconds.",
		"random":  getRandomNumber(),
		"headers": json{
			"expires": expires,
		},
	}
}

// _pragmaNoCache: Creates a JSON with non cacheable header using the `Pragma`
// response header. This header is pre HTTP-1.1 and is deprecated. Exists only
// for backward compatibility.
func _pragmaNoCache(ctx *gin.Context) (int, json) {
	pragma := "no-cache"
	ctx.Header("Pragma", pragma)

	return http.StatusOK, json{
		"message": "Using the `Pragma` header to set caching as disabled.You can see the `random` value change every time you make a request",
		"random":  getRandomNumber(),
		"headers": json{
			"pragma": pragma,
		},
	}
}

func getTime(addSec int) string {
	now := time.Now().Add(time.Second * time.Duration(addSec)).UTC()
	return now.Format(time.RFC1123)
}

func getRandomNumber() int {
	randInt := rand.Intn(200)
	return randInt
}
