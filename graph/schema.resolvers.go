package graph

import (
	"context"
	"fmt"
	"log"
	"the-rush/graph/model"
	"the-rush/graph/runtime"
)

func (r *mutationResolver) CreatePlayer(ctx context.Context, playerInput model.PlayerInput) (*model.Player, error) {
	log.Println("mutationResolver.CreatePlayer")
	player, err := playerInput.ToPlayer()
	if err != nil {
		return nil, fmt.Errorf("failed to convert from player input: %w", err)
	}
	return r.repo.Player().Create(ctx, player)
}

func (r *queryResolver) Player(ctx context.Context, id string) (*model.Player, error) {
	log.Println("queryResolver.Player")
	return r.repo.Player().Read(ctx, id)
}

func (r *queryResolver) Players(ctx context.Context, args model.PlayersArgs) (*model.PlayersResponse, error) {
	log.Println("queryResolver.Players")
	pagination, err := model.PlayerPaginationFromArgs(args)
	if err != nil {
		return new(model.PlayersResponse), fmt.Errorf("failed to parse players args: %w", err)
	}

	resp, err := r.repo.Player().Search(ctx, args.Name, *pagination, false)
	if err != nil {
		return new(model.PlayersResponse), fmt.Errorf("failed to parse players args: %w", err)
	}
	if resp == nil {
		return new(model.PlayersResponse), nil
	}
	return resp, nil
}

// Mutation returns runtime.MutationResolver implementation.
func (r *Resolver) Mutation() runtime.MutationResolver { return &mutationResolver{r} }

// Query returns runtime.QueryResolver implementation.
func (r *Resolver) Query() runtime.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
