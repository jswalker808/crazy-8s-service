package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayWebsocketProxyRequest) (interface{}, error) {

	log.Println("In the handler func")

	switch routeKey := request.RequestContext.RouteKey; routeKey {
		case "$connect": 
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body: "Successfully connected",
			}, nil
		case "$disconnect":
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body: "Successfully disconnected",
			}, nil
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body: "Play the game",
			}, nil
	}
}

func main() {
	lambda.Start(handler)
}
