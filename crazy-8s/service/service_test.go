// +build integration

package service

import (
	"context"
	"crazy-8s/config"
	"crazy-8s/notification"
	"crazy-8s/repository"
	"crazy-8s/testutil"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/suite"
)

type GameServiceTestSuite struct {
	suite.Suite
	dynamoDbContainer *testutil.DynamoDbContainer
	service  *GameService
	ctx context.Context
}

type MockNotifier struct {}

func (mock *MockNotifier) Send(connectionId string, bytes []byte) error { return nil }
func (mock *MockNotifier) SendAll(connectionMap map[string][]byte) []notification.NotificationError { return make([]notification.NotificationError, 0) }


func (suite *GameServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	dynamoDbContainer, err := testutil.CreateDynamoDbContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.dynamoDbContainer = dynamoDbContainer

	repository := repository.NewGameRepository(dynamodb.NewFromConfig(config.LoadAwsConfig()))
	suite.service = NewGameService(repository, &MockNotifier{})
}

func (suite *GameServiceTestSuite) TearDownSuite() {
	if err := suite.dynamoDbContainer.Container.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating dynamodb container: %s", err)
	}
}


func (suite *GameServiceTestSuite) TestCreateGame() {

}

func (suite *GameServiceTestSuite) TestJoinGame() {

}

func TestIntegrationGameServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GameServiceTestSuite))
}
