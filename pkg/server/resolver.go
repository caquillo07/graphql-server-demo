package server

import (
	"context"

	"github.com/caquillo07/graphql-server-demo/pkg/gqlgen/schema"
	gqlServer "github.com/caquillo07/graphql-server-demo/pkg/gqlgen/server"
)

func (s *server) Mutation() gqlServer.MutationResolver {
	return s
}

func (s *server) Query() gqlServer.QueryResolver {
	return s
}

func (s *server) UpdateUser(ctx context.Context, name string) (*schema.User, error) {
	return &schema.User{
		ID:   "123e4567-e89b-12d3-a456-426655440000",
		Name: name,
	}, nil
}

func (s *server) GetUsers(ctx context.Context) ([]*schema.User, error) {
	return []*schema.User{
		{
			ID:   "123e4567-e89b-12d3-a456-426655440000",
			Name: "Bob",
		},
	}, nil
}
