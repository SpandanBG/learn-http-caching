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

	go _itemModify()

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
	r.GET("/cache_control_5_sec", jsonRouteWrapper(_cacheControl5Sec))
	r.GET("/cache_control_no_store", jsonRouteWrapper(_cacheControlNoStore))
	r.GET("/cache_control_no_cache", jsonRouteWrapper(_cacheControlNoCache))
	r.GET("/cache_control_must_revalidate", jsonRouteWrapper(_cacheControlMustRevalidate))
	r.GET("/cache_control_last_modified", jsonRouteWrapper(_cacheControlLastModfied))
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

// _cacheControl5Sec: Creates a JSON with a `cache-content: max-age=5` response
// header.
func _cacheControl5Sec(ctx *gin.Context) (int, json) {
	ctx.Header("Cache-Control", "max-age=5")

	return http.StatusOK, json{
		"message": "Using the `cache-control` header to set caching for 5 seconds. You can see the `random` value change if you request for it after 5 seconds.",
		"random":  getRandomNumber(),
		"headers": json{
			"Cache-Control": "max-age=5",
		},
	}
}

// _cacheControlNoStore: Creates a JSON with `cache-control: no-store` response
// header which make the client not cache the content
func _cacheControlNoStore(ctx *gin.Context) (int, json) {
	ctx.Header("Cache-Control", "no-store")

	return http.StatusOK, json{
		"message": "Using the `cache-control` header to ensure no caching. You can see the `random` value change everytime.",
		"random":  getRandomNumber(),
		"headers": json{
			"cache-control": "no-store",
		},
	}
}

var _cacheControlNoCacheData string

// _cacheControlNoCache: Creates a JSON with `cache-control: no-cache` and an
// `ETag: ${Token}` response header. The response will be different if the
// `Token` query param is different from the last passed value.
func _cacheControlNoCache(ctx *gin.Context) (int, json) {
	token := ctx.Query("Token")

	if token == _cacheControlNoCacheData {
		return http.StatusNotModified, nil
	}

	_cacheControlNoCacheData = token
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("ETag", token)

	return http.StatusOK, json{
		"message": "Using the `cache-control` and `etag` headers to ensure cached value is revalidated with `ETag`. You can see the `random` value change if the `Token` changes.",
		"token":   token,
		"random":  getRandomNumber(),
		"headers": json{
			"Cache-Control": "no-cache",
			"ETag":          token,
		},
	}
}

var _cacheControlMustRevalidateData string

// _cacheControlMustRevalidate: Creates a JSON with `cache-control: max-age=5,
// must-revalidate, private` and an `ETag: ${Token}`response header. When client
// calls with no network, the cache will not be served to the user. The response
// will be different if the `Token` query param is different from the last passed
// value.
func _cacheControlMustRevalidate(ctx *gin.Context) (int, json) {
	token := ctx.Query("Token")

	if token == _cacheControlMustRevalidateData {
		return http.StatusNotModified, nil
	}

	_cacheControlMustRevalidateData = token

	ctx.Header("Cache-Control", "max-age=5, must-revalidate, private")
	ctx.Header("ETag", token)

	return http.StatusOK, json{
		"message": "Using the `cache-control` and `etag` headers to ensure cached value is revalidated after value is stale and checks `ETag`. This ensures to revalidate when the client has nextwork after the cache is stale. You can see the `random` value change if `Token` changes.",
		"token":   token,
		"random":  getRandomNumber(),
		"headers": json{
			"Cache-Control": "max-age=5, must-revalidate, private",
			"ETag":          token,
		},
	}
}

// _cacheControlLastModfied: Cretaes a JSON with with `cache-control: max-age=5,
// must-revalidate, private` and an `Last-Modified: XYZ` response headers. When
// the client requests with the `if-modified-since` header, it checks with the
// Item's last modified time. If it is different, it will return the new data,
// otherwise return 304. The Item will keep changing every 10 seconds since the
// start of the server
func _cacheControlLastModfied(ctx *gin.Context) (int, json) {
	lastModified := ctx.GetHeader("if-modified-since")
	if lastModified == _itemLastModified {
		return http.StatusNotModified, nil
	}

	ctx.Header("Cache-Control", "max-age=5, must-revalidate, private")
	ctx.Header("Last-Modified", _itemLastModified)

	return http.StatusOK, json{
		"message": "Using the `cache-control` and `last-modified` response header for caching. The random number will keep changing after 10 seconds since the start of the server. However the client will only revalidate after 5 seconds when local cache is stale.",
		"random":  _item,
		"headers": json{
			"Cache-Control": "max-age=5, must-revalidate, private",
			"Last-Modified": _itemLastModified,
		},
	}
}

var _item int
var _itemLastModified string

func _itemModify() {
	for {
		time.Sleep(10 * time.Second)
		_item = getRandomNumber()
		_itemLastModified = getTime(0)
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
