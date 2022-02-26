package session

import (
	"context"
	"fmt"

	"api-backend-template/internal/constants"
	"api-backend-template/internal/utils"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

var (
	Region  = utils.LoadEnv("Region", "ap-southeast-1")
	EnvType = utils.LoadEnv("EnvType", "dev")

	DdbTableName = fmt.Sprintf(constants.DB_FMT, constants.ServiceName, EnvType) // api-backend-template-dev

	// Using the Config value, create the DynamoDB client
	Ddbsvc = getDdbClient()
)

func getDdbClient() *dynamodb.Client {
	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(Region),
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), loadOpts...)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return dynamodb.NewFromConfig(cfg)
}
