package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin/v2"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"golang.org/x/net/context/ctxhttp"
)

var tracingClient = apmhttp.WrapClient(http.DefaultClient)

func main() {
	// os.Setenv("ELASTIC_APM_SERVER_URL", "http://apm-server:8200")
	// os.Setenv("ELASTIC_APM_SERVICE_NAME", "ms-payments")

	fmt.Println("Start ms-payments")

	r := gin.Default()
	r.Use(apmgin.Middleware(r))

	// For healthcheck
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	r.GET("/payments", func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "payments.handler", "app")
		defer span.End()

		payments, err := getPayments(ctx)
		if err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"payments": payments,
		})
	})
	r.Run()
}

func getPayments(ctx context.Context) (map[string]interface{}, error) {

	span, ctx := apm.StartSpan(ctx, "getTodoFromAPI", "app")
	defer span.End()

	var result map[string]interface{}

	resp, err := ctxhttp.Get(ctx, tracingClient, "https://jsonplaceholder.typicode.com/todos/1")
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
