package spannerclient

import (
	"context"
	"time"

	"cloud.google.com/go/spanner"
)

// DataClientOption is option that can be applied when calling NewDataClient.
type DataClientOption func(c *dataClientOptions)

// NewDataClient provides a spanner.Client.
func NewDataClient(dbName string, opts ...DataClientOption) (*spanner.Client, error) {
	options := &dataClientOptions{
		databaseName: dbName,

		// default values for optional options
		createClientTimeout:           time.Minute,
		fractionPreparedWriteSessions: 1.0,
		maxIdleSessions:               10,
		minOpenedSessions:             10,
		maxBurstSessions:              50,
	}
	// apply optional options
	for _, opt := range opts {
		opt(options)
	}
	return newDataClient(options)
}

// CreateClientTimeoutOption sets timeout when creating Spanner data client
var CreateClientTimeoutOption = func(timeout time.Duration) DataClientOption {
	return func(c *dataClientOptions) {
		c.createClientTimeout = timeout
	}
}

type dataClientOptions struct {
	databaseName string

	// session pool
	createClientTimeout           time.Duration
	fractionPreparedWriteSessions float64
	maxIdleSessions               uint64
	minOpenedSessions             uint64
	maxBurstSessions              uint64
}

func newDataClient(opts *dataClientOptions) (*spanner.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opts.createClientTimeout)
	defer cancel()

	cl, err := spanner.NewClientWithConfig(ctx, opts.databaseName, spanner.ClientConfig{
		NumChannels: 10,
		SessionPoolConfig: spanner.SessionPoolConfig{
			WriteSessions: opts.fractionPreparedWriteSessions,
			MaxIdle:       opts.maxIdleSessions,
			MinOpened:     opts.minOpenedSessions,
			MaxBurst:      opts.maxBurstSessions,
		},
	})
	if err != nil {
		return nil, err
	}
	return cl, nil
}
