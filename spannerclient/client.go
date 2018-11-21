package spannerclient

import (
	"time"
	"cloud.google.com/go/spanner/admin/database/apiv1"
	"context"
	"log"
	"cloud.google.com/go/spanner"
)

func NewAdminClient() *database.DatabaseAdminClient {
	ctx, cancel := context.WithTimeout(context.Background(), CreateSchemaTimeout)
	defer cancel()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	return adminClient
}

func NewDataClient() (*spanner.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute * 1)
	defer cancel()

	cl, err := spanner.NewClientWithConfig(ctx, DatabaseName, spanner.ClientConfig{
		NumChannels: 10,
		SessionPoolConfig: spanner.SessionPoolConfig{
			MaxBurst: 50,
		},
	})
	if err != nil {
		return nil, err
	}
	return cl, nil
}
