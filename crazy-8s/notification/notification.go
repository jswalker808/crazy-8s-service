package notification

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

type Notifier interface {
	Send(string, []byte)
}

type ApiGatewayNotifier struct {
	apiGatewayClient *apigatewaymanagementapi.Client
}

func NewApiGatewayNotifier() *ApiGatewayNotifier {
	return &ApiGatewayNotifier{}
}

func (apiGatewayNotifier *ApiGatewayNotifier) SetClient(apiGatewayClient *apigatewaymanagementapi.Client) {
	apiGatewayNotifier.apiGatewayClient = apiGatewayClient
}

func (apiGatewayNotifier *ApiGatewayNotifier) Send(connectionId string, bytes []byte) error {
	log.Println(connectionId)

	input := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionId),
		Data: bytes,
	}
	
	_, err := apiGatewayNotifier.apiGatewayClient.PostToConnection(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}



