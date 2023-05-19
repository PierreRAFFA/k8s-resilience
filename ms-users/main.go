package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context/ctxhttp"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin/v2"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
)

var tracingClient = apmhttp.WrapClient(http.DefaultClient)

func main() {
	// os.Setenv("ELASTIC_APM_SERVER_URL", "http://apm-server:8200")
	// os.Setenv("ELASTIC_APM_SERVICE_NAME", "ms-users")

	fmt.Println("Start ms-users")

	r := gin.Default()
	r.Use(apmgin.Middleware(r))

	// For healthcheck
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	r.GET("/friends", func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "users_friends.handler", "app")
		defer span.End()

		friends, err := getUserFriends(ctx)
		if err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"friends": friends,
		})
	})
	r.Run()
}

// ms-friends does not exist and generates an error in ElasticAPM
func getUserFriends(ctx context.Context) (map[string]interface{}, error) {
	span, ctx := apm.StartSpan(ctx, "getUserFriends", "app")
	defer span.End()

	var result map[string]interface{}

	resp, err := ctxhttp.Get(ctx, tracingClient, "https://jsonplaceholder.typicode.com/todos/2")
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return result, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, err
}
