package notification

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

type Notifier interface {
	Send(string, []byte) NotificationError
	SendAll(map[string][]byte) []NotificationError
}

type NotificationError struct {
	ConnectionId string
	Error error
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

func (apiGatewayNotifier *ApiGatewayNotifier) Send(connectionId string, bytes []byte) NotificationError {
	log.Println(connectionId)

	input := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionId),
		Data: bytes,
	}
	
	_, err := apiGatewayNotifier.apiGatewayClient.PostToConnection(context.TODO(), input)
	if err != nil {
		return NotificationError{connectionId, err}
	}

	return NotificationError{}
}

func (apiGatewayNotifier *ApiGatewayNotifier) SendAll(notificationMap map[string][]byte) []NotificationError {
	numNotifications := len(notificationMap)
	resultChannel := make(chan NotificationError, numNotifications)

	for connectionId, bytes := range notificationMap {
		go func(connectionId string, bytes []byte) {
			resultChannel <- apiGatewayNotifier.Send(connectionId, bytes)
		}(connectionId, bytes)
	}

	notificationErrors := make([]NotificationError, 0, numNotifications)
	for i := 0; i < numNotifications; i++ {
		if notificationError := <-resultChannel; notificationError != (NotificationError{}) {
			log.Printf("Connection %v ran into an error: %v", notificationError.ConnectionId, notificationError.Error)
			notificationErrors = append(notificationErrors, notificationError)
		}
	}

	close(resultChannel)

	return notificationErrors
}


