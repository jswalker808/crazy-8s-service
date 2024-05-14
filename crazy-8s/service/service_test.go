package service

import (
	"context"
	"crazy-8s/notification"
	"crazy-8s/repository"
	"crazy-8s/testutil"
	"crazy-8s/transport"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MockNotifier struct {}
func (mock *MockNotifier) Send(connectionId string, bytes []byte) error { return nil }
func (mock *MockNotifier) SendAll(connectionMap map[string][]byte) []notification.NotificationError { return make([]notification.NotificationError, 0) }

type GameServiceTestSuite struct {
	suite.Suite
	dynamoDbContainer *testutil.DynamoDbContainer
	service  *GameService
	ctx context.Context
}

func (suite *GameServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	dynamoDbContainer, err := testutil.CreateDynamoDbContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				URL:           dynamoDbContainer.Endpoint,
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})
	
	config, err := awsConfig.LoadDefaultConfig(suite.ctx, awsConfig.WithRegion("us-west-1"), awsConfig.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatal(err)
	}

	tableName := "TestTable"
	os.Setenv("TABLE_NAME", tableName)

	dynamoDbClient := dynamodb.NewFromConfig(config)
	_, createTableErr := dynamoDbClient.CreateTable(suite.ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("gameId"),
				KeyType:       types.KeyTypeHash,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("gameId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if (createTableErr != nil) {
		log.Fatal(err)
	}

	repository := repository.NewGameRepository(dynamoDbClient)
	suite.service = NewGameService(repository, &MockNotifier{})
}

func (suite *GameServiceTestSuite) TestCreateGame() {
	if (testing.Short()) {
		suite.T().Skip("skipping integration test")
	}

	assert := assert.New(suite.T())

	connectionId := "testConnectionId1"
	createGameRequest := transport.CreateGameRequest{
		PlayerName: "player1",
	}
	createdGame, err := suite.service.CreateGame(connectionId, &createGameRequest)
	assert.NoError(err)
	assert.NotNil(createdGame)
}

func (suite *GameServiceTestSuite) TestJoinGame() {
	if (testing.Short()) {
		suite.T().Skip("skipping integration test")
	}

	assert := assert.New(suite.T())

	connectionId := "testConnectionId1"
	createGameRequest := transport.CreateGameRequest{
		PlayerName: "player1",
	}
	createdGame, err := suite.service.CreateGame(connectionId, &createGameRequest)
	assert.NoError(err)

	joinGameRequest := transport.JoinGameRequest{
		PlayerName: "player2",
		GameId: createdGame.GetId(),
	}
	updatedGame, err := suite.service.JoinGame(connectionId, &joinGameRequest)
	assert.NoError(err)
	assert.NotNil(updatedGame)
	assert.Equal(2, len(updatedGame.GetPlayers()))
}

func TestIntegrationGameServiceTestSuite(t *testing.T) {
	if (testing.Short()) {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(GameServiceTestSuite))
}
