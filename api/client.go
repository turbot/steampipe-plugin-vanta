package api

import (
	"context"

	"github.com/machinebox/graphql"
)

// Vanta API Client
type Client struct {
	SessionId *string
	Token     *string
	Graphql   *graphql.Client
}

func CreateClient(ctx context.Context, config ClientConfig) (*Client, error) {
	return &Client{
		Token:   config.ApiToken,
		Graphql: graphql.NewClient("https://api.vanta.com/graphql"),
	}, nil
}

func CreateAppClient(ctx context.Context, config ClientConfig) (*Client, error) {
	return &Client{
		SessionId: config.SessionId,
		Graphql:   graphql.NewClient("https://app.vanta.com/graphql"),
	}, nil
}
