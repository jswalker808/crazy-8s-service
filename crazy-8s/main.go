package main

import (
	"context"
	"crazy-8s/api"
	"crazy-8s/config"
	"crazy-8s/global"
	"crazy-8s/notification"
	"crazy-8s/repository"
	"crazy-8s/service"
	"crazy-8s/transport"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type APIGatewayRequest = events.APIGatewayWebsocketProxyRequest
type APIGatewayResponse = events.APIGatewayProxyResponse

var router *api.Router
var awsConfig aws.Config

func handler(ctx context.Context, request APIGatewayRequest) (APIGatewayResponse, error) {

	log.Println("In the handler function")

	ctx = context.WithValue(ctx, global.ConnectionIdCtxKey{}, request.RequestContext.ConnectionID)

	switch routeKey := request.RequestContext.RouteKey; routeKey {
	case "$connect":
		return APIGatewayResponse{
			StatusCode: 200,
			Body:       "Successfully connected",
		}, nil
	case "$disconnect":
		return APIGatewayResponse{
			StatusCode: 200,
			Body:       "Successfully disconnected",
		}, nil
	default:
		return handleGamePlay(ctx, request)
	}
}

func handleGamePlay(ctx context.Context, apiGatewayRequest APIGatewayRequest) (APIGatewayResponse, error) {
	log.Printf("Request: %v", apiGatewayRequest.Body)

	baseRequest, err := transport.NewBaseRequest(apiGatewayRequest.Body)
	if err != nil {
		return APIGatewayResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, err
	}

	gameRequest, err := transport.NewGameRequest(baseRequest)
	if err != nil {
		return APIGatewayResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, err
	}

	apiGatewayClient := loadApiGatewayClient(awsConfig, apiGatewayRequest.RequestContext.DomainName, apiGatewayRequest.RequestContext.Stage)
	router.GameService().Notifier().(*notification.ApiGatewayNotifier).SetClient(apiGatewayClient)

	if err := router.HandleRequest(ctx, baseRequest.Action, gameRequest); err != nil {
		return APIGatewayResponse{
			StatusCode: 400,
		}, err
	}

	return APIGatewayResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	awsConfig = config.LoadAwsConfig()

	dynamoDbClient := dynamodb.NewFromConfig(awsConfig)
	gameRepository := repository.NewGameRepository(dynamoDbClient)
	connectionRepository := repository.NewConnectionRepository(dynamoDbClient)

	apiGatewayNotifier := notification.NewApiGatewayNotifier()

	gameService := service.NewGameService(gameRepository, connectionRepository, apiGatewayNotifier)

	router = api.NewRouter(gameService)

	lambda.Start(handler)
}

func loadApiGatewayClient(awsConfig aws.Config, domain string, stage string) *apigatewaymanagementapi.Client {
	endpoint := fmt.Sprintf("https://%v/%v", domain, stage)
	return apigatewaymanagementapi.NewFromConfig(awsConfig.Copy(), func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
}
