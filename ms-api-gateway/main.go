package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context/ctxhttp"

	"os"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin/v2"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
)

var tracingClient = apmhttp.WrapClient(http.DefaultClient)

var usersServiceUrl = os.Getenv("USERS_SERVICE_URL")
var paymentsServiceUrl = os.Getenv("PAYMENTS_SERVICE_URL")

func main() {

	//os.Setenv("ELASTIC_APM_SERVER_URL", "http://apm-server:8200")
	//os.Setenv("ELASTIC_APM_SERVICE_NAME", "ms-users")

	r := gin.Default()
	r.Use(apmgin.Middleware(r))

	// For healthcheck
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	r.GET("/api/users/payments", func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "users_payments.handler", "app")
		defer span.End()

		res, err := getUserPayments(ctx)
		if err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"payments": res,
		})
	})

	r.GET("/api/users/friends", func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "users_friends.handler", "app")
		defer span.End()

		res, err := getUserFriends(ctx)
		if err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"friends": res,
		})
	})
	r.Run()
}

func getUserPayments(ctx context.Context) (map[string]interface{}, error) {
	span, ctx := apm.StartSpan(ctx, "getUserPayments", "app")
	defer span.End()

	var result map[string]interface{}

	resp, err := ctxhttp.Get(ctx, tracingClient, fmt.Sprintf("%s%s", paymentsServiceUrl, "/api/payments"))
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

// ms-friends does not exist and generates an error in ElasticAPM
func getUserFriends(ctx context.Context) (map[string]interface{}, error) {
	span, ctx := apm.StartSpan(ctx, "getUserFriends", "app")
	defer span.End()

	var result map[string]interface{}

	resp, err := ctxhttp.Get(ctx, tracingClient, fmt.Sprintf("%s%s", usersServiceUrl, "/api/users/friends"))
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
