package service

import (
	"context"
	repositoryPkg "crazy-8s/repository"
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

type GameServiceTestSuite struct {
	suite.Suite
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
	connectionTableName := "ConnectionsTable"
	os.Setenv("TABLE_NAME", tableName)
	os.Setenv("CONNECTIONS_TABLE_NAME", connectionTableName)

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

	_, createConnTableErr := dynamoDbClient.CreateTable(suite.ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(connectionTableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("connectionId"),
				KeyType:       types.KeyTypeHash,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("connectionId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if (createConnTableErr != nil) {
		log.Fatal(err)
	}

	gameRepository := repositoryPkg.NewGameRepository(dynamoDbClient)

	suite.service = NewGameService(gameRepository)
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

func (suite *GameServiceTestSuite) TestLeaveGame() {
	if (testing.Short()) {
		suite.T().Skip("skipping integration test")
	}

	assert := assert.New(suite.T())

	ownerConnectionId := "testConnectionId3"
	secondPlayerConnectionId := "testConnectionId4"

	createGameRequest := transport.CreateGameRequest{
		PlayerName: "player3",
	}
	createdGame, err := suite.service.CreateGame(ownerConnectionId, &createGameRequest)
	assert.NoError(err)

	joinGameRequest := transport.JoinGameRequest{
		PlayerName: "player4",
		GameId: createdGame.GetId(),
	}
	updatedGame, err := suite.service.JoinGame(secondPlayerConnectionId, &joinGameRequest)
	assert.NoError(err)
	assert.NotNil(updatedGame)
	assert.Equal(2, len(updatedGame.GetPlayers()))
	
	updatedGame, err = suite.service.LeaveGame(secondPlayerConnectionId)
	assert.NoError(err)
	assert.NotNil(updatedGame)
	assert.Equal(1, len(updatedGame.GetPlayers()))
}

func TestIntegrationGameServiceTestSuite(t *testing.T) {
	if (testing.Short()) {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(GameServiceTestSuite))
}
