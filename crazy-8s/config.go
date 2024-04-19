package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

func loadAwsConfig() aws.Config {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-1"))
	if err != nil {
		log.Panic("Unable to load aws configuration")
	}
	return config
}

func loadApiGatewayClient(awsConfig aws.Config, domain string, stage string) *apigatewaymanagementapi.Client {
	endpoint := fmt.Sprintf("https://%v/%v", domain, stage)
	return apigatewaymanagementapi.NewFromConfig(awsConfig.Copy(), func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
}