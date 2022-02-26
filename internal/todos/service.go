package todos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api-backend-template/internal/constants"
	"api-backend-template/internal/session"
	"api-backend-template/internal/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

var (
	defaultPageSize int32 = 20
)

//go:generate mockgen -destination=service_mock.go -package=todo . Service

// Service ...
type Service interface {
	GetItems(*GetItemInput) (*GetItemOutput, error)
	GetItem(id string) (*Todo, error)
	CreateItem(item Todo) (*Todo, error)
	UpdateItem(item Todo) (*Todo, error)
	DeleteItem(id string) error
}

type DdbClient struct {
	DB        *dynamodb.Client
	TableName string
	PageSize  int32
}

func NewWithDefaults() Service {
	return &DdbClient{
		DB:        session.Ddbsvc,
		TableName: session.DdbTableName,
		PageSize:  defaultPageSize,
	}
}

func NewService(db *dynamodb.Client, tableName string, pageSize int32) Service {
	return &DdbClient{
		DB:        db,
		TableName: tableName,
		PageSize:  pageSize,
	}
}

type GetItemInput struct {
	LastKey *string
}

type GetItemOutput struct {
	Items   []Todo
	LastKey *string
	Count   int32
}

func (ddbclient *DdbClient) GetItems(input *GetItemInput) (*GetItemOutput, error) {
	qInput := &dynamodb.QueryInput{
		TableName: aws.String(ddbclient.TableName),
		ExpressionAttributeNames: map[string]string{
			"#pk": "pk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk_val": &types.AttributeValueMemberS{Value: constants.Todo},
		},
		KeyConditionExpression: aws.String("#pk = :pk_val"),
		Limit:                  aws.Int32(ddbclient.PageSize),
	}
	if input.LastKey != nil {
		qInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: constants.Todo},
			"sk": &types.AttributeValueMemberS{Value: *input.LastKey},
		}
	}

	result, err := ddbclient.DB.Query(context.TODO(), qInput)
	if err != nil {
		return nil, err
	}

	Todos := make([]Todo, 0)
	attributevalue.UnmarshalListOfMaps(result.Items, &Todos)

	output := &GetItemOutput{
		Items: Todos,
		Count: result.Count,
	}
	if _, ok := result.LastEvaluatedKey["sk"]; ok {
		output.LastKey = utils.String(result.LastEvaluatedKey["sk"].(*types.AttributeValueMemberS).Value)
	}
	return output, nil
}

func (ddbclient *DdbClient) GetItem(id string) (*Todo, error) {
	item, err := ddbclient.DB.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(ddbclient.TableName),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: constants.Todo},
			"sk": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}

	if item.Item == nil {
		return nil, nil
	}

	Todo := &Todo{}
	err = attributevalue.UnmarshalMap(item.Item, Todo)
	return Todo, err
}

func (ddbclient *DdbClient) GetItemByQueryExecutionID(qID string) (*Todo, error) {
	qInput := &dynamodb.QueryInput{
		TableName: aws.String(ddbclient.TableName),
		ExpressionAttributeNames: map[string]string{
			"#pk":  "pk",
			"#qid": "athena_query_execution_id",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk_val":  &types.AttributeValueMemberS{Value: constants.Todo},
			":qid_val": &types.AttributeValueMemberS{Value: qID},
		},
		KeyConditionExpression: aws.String("#pk = :pk_val"),
		FilterExpression:       aws.String("#qid = :qid_val"),
	}

	result, err := ddbclient.DB.Query(context.TODO(), qInput)
	if err != nil {
		return nil, err
	}

	Todos := make([]Todo, 0)
	attributevalue.UnmarshalListOfMaps(result.Items, &Todos)

	if len(Todos) == 0 {
		return nil, nil
	}

	if len(Todos) != 1 {
		return nil, fmt.Errorf("more than one item found for given qID (%s)", qID)
	}

	return &Todos[0], nil
}

func (ddbclient *DdbClient) CreateItem(item Todo) (*Todo, error) {
	if item.ID == uuid.Nil {
		id, _ := uuid.NewRandom()
		item.ID = id
	}

	return ddbclient.UpdateItem(item)
}

func (ddbclient *DdbClient) UpdateItem(item Todo) (*Todo, error) {
	t := time.Now().Format(time.RFC3339)
	item.PartitionKey = constants.Todo
	item.SortKey = item.ID.String()
	if item.CreatedAt == "" {
		item.CreatedAt = t
	}
	item.UpdatedAt = t

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, errors.New("Unable to getPutItemInput: " + err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(ddbclient.TableName),
	}

	_, err = ddbclient.DB.PutItem(context.TODO(), input)
	return &item, err
}

func (ddbclient *DdbClient) DeleteItem(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(ddbclient.TableName),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: constants.Todo},
			"sk": &types.AttributeValueMemberS{Value: id},
		},
	}
	_, err := ddbclient.DB.DeleteItem(context.TODO(), input)
	return err
}
