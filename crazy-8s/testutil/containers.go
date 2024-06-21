package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DynamoDbContainer struct {
	Container testcontainers.Container
	Endpoint string
}

func CreateDynamoDbContainer(ctx context.Context) (*DynamoDbContainer, error) {
    req := testcontainers.ContainerRequest{
        Image:        "amazon/dynamodb-local",
        ExposedPorts: []string{"8000/tcp"},
        WaitingFor:   wait.ForLog("Initializing DynamoDB Local").WithStartupTimeout(time.Minute * 2),
    }

    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        return nil, err
    }

	host, err := container.Host(ctx)
    if err != nil {
        return nil, err
    }

    port, err := container.MappedPort(ctx, "8000")
    if err != nil {
        return nil, err
    }

    endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())

    return &DynamoDbContainer{
		Container: container,
		Endpoint: endpoint,
	}, nil
}