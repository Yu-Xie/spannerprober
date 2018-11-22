package spannerclient

import (
	"context"
	"time"

	"cloud.google.com/go/spanner/admin/database/apiv1"
)

// AdminClientOption is option that can be applied when calling NewAdminClient.
type AdminClientOption func(*adminClientOptions)

// NewAdminClient provides a Spanner admin client.
func NewAdminClient(opts ...AdminClientOption) (*database.DatabaseAdminClient, error) {
	options := &adminClientOptions{
		createClientTimeout: time.Minute,
	}
	for _, opt := range opts {
		opt(options)
	}
	return newAdminClient(options)
}

type adminClientOptions struct {
	createClientTimeout time.Duration
}

func newAdminClient(opts *adminClientOptions) (*database.DatabaseAdminClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opts.createClientTimeout)
	defer cancel()
	return database.NewDatabaseAdminClient(ctx)
}
