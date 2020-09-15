package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"the-rush/constant"
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

func (r *queryResolver) TotalYardsByTeam(ctx context.Context, order model.Order) ([]*model.TotalYardsByTeamResponse, error) {
	log.Println("queryResolver.TotalYardsByTeam")
	playersMap, err := r.repo.HGetAll(ctx, constant.RedisKeyPlayerRecords.String()).Result()
	if err != nil {
		return nil, err
	}

	players := make([]model.Player, 0)
	for _, v := range playersMap {
		player := new(model.Player)
		err = json.Unmarshal([]byte(v), player)
		if err != nil {
			return nil, err
		}
		players = append(players, *player)
	}

	totalYardsByTeam := make(map[string]int)
	for _, v := range players {
		totalYardsByTeam[v.Team] += v.TotalRushingYards
	}

	resp := make([]*model.TotalYardsByTeamResponse, 0)
	for k, v := range totalYardsByTeam {
		r := &model.TotalYardsByTeamResponse{
			Team:       k,
			TotalYards: v,
		}
		resp = append(resp, r)
	}

	sortable := model.SortableTotalYardsByTeamResponse(resp)
	if order == model.OrderAsc {
		sort.Sort(sortable)
	} else {
		sort.Sort(sort.Reverse(sortable))
	}

	return sortable, nil
}

// Mutation returns runtime.MutationResolver implementation.
func (r *Resolver) Mutation() runtime.MutationResolver { return &mutationResolver{r} }

// Query returns runtime.QueryResolver implementation.
func (r *Resolver) Query() runtime.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
