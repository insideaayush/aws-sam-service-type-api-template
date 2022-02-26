package main

import (
	"context"
	"os"

	"api-backend-template/internal/log"
	"api-backend-template/internal/routes"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var ginLambda *ginadapter.GinLambda

// Handler is the main entry point for Lambda. Receives a proxy request and
// returns a proxy response
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Init(req.HTTPMethod, req.Path)

	if ginLambda == nil {
		r := gin.Default()
		r.Use(log.Logger(logrus.StandardLogger()), gin.Recovery()) // Add this middleware to use logrus with gin

		// Map routes to handlers
		routes.RegisterRoutes(r.Group("/api"))

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func init() {
	if os.Getenv("EnvType") != "dev" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	lambda.Start(Handler)
}
