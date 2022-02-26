require (
	github.com/aws/aws-lambda-go v1.23.0
	github.com/aws/aws-sdk-go-v2 v1.13.0
	github.com/aws/aws-sdk-go-v2/config v1.13.1
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.6.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.13.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.12.0
	github.com/gin-gonic/gin v1.7.7
	github.com/google/uuid v1.3.0
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8

module api-backend-template

go 1.16
