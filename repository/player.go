package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
	uuid2 "github.com/google/uuid"

	"the-rush/graph/model"
)

type PlayerRepository interface {
	Create(ctx context.Context, player *model.Player) (*model.Player, error)
	Read(ctx context.Context, id string) (*model.Player, error)
	Search(ctx context.Context, name *string, pagination model.PlayerPagination, withoutLimit bool) (*model.PlayersResponse, error)
}

type playerRepository struct{ *Repository }

func (x *playerRepository) Create(ctx context.Context, player *model.Player) (*model.Player, error) {
	playerJson, err := json.Marshal(player)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}
	err = x.Client.Watch(ctx, func(tx *redis.Tx) error {
		var e error
		_, e = tx.HSet(ctx, x.playerRecords, player.ID, string(playerJson)).Result()
		if e != nil {
			return e
		}
		_, e = tx.ZAdd(ctx, x.playerByCreatedBy, &redis.Z{Score: float64(player.CreatedAt.UnixNano()), Member: player.ID}).Result()
		if e != nil {
			return e
		}
		_, e = tx.ZAdd(ctx, x.playerByLongestRush, &redis.Z{Score: float64(player.LongestRush.Value), Member: player.ID}).Result()
		if e != nil {
			return e
		}
		_, e = tx.ZAdd(ctx, x.playerByTotalRushingTouchdowns, &redis.Z{Score: float64(player.TotalRushingTouchdowns), Member: player.ID}).Result()
		if e != nil {
			return e
		}
		_, e = tx.ZAdd(ctx, x.playerByTotalRushingYards, &redis.Z{Score: float64(player.TotalRushingYards), Member: player.ID}).Result()
		if e != nil {
			return e
		}
		names := strings.Split(player.Name, " ")
		for _, name := range names {
			_, e = tx.ZAdd(ctx, fmt.Sprintf(x.playerByNamePrefix, strings.ToLower(name)), &redis.Z{Score: 1.0, Member: player.ID}).Result()
			if e != nil {
				return e
			}
		}
		return nil
	}, []string{x.playerRecords, x.playerByNamePrefix, x.playerByCreatedBy, x.playerByLongestRush, x.playerByTotalRushingTouchdowns, x.playerByTotalRushingYards}...)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return player, nil
}

func (x *playerRepository) Read(ctx context.Context, id string) (*model.Player, error) {
	str, err := x.Client.HGet(ctx, x.playerRecords, id).Result()
	if err != nil && strings.Contains(err.Error(), "redis: nil") {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to hget: %w", err)
	}

	player := new(model.Player)
	err = json.Unmarshal([]byte(str), player)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return player, nil
}

func (x *playerRepository) Search(
	ctx context.Context,
	name *string,
	pagination model.PlayerPagination,
	withoutLimit bool,
) (*model.PlayersResponse, error) {
	if name != nil && *name != "" {
		return x.searchByNameAndPagination(ctx, strings.ToLower(*name), pagination, withoutLimit)
	}
	return x.searchByIndexAndPagination(ctx, pagination.OrderBy, pagination, withoutLimit)
}

func (x *playerRepository) searchByNameAndPagination(
	ctx context.Context,
	nameStr string,
	pagination model.PlayerPagination,
	withoutLimit bool,
) (*model.PlayersResponse, error) {
	var (
		names         = strings.Split(nameStr, " ")
		interStoreIDs = make([]string, 0)
		resultID      string
	)

	err := x.Watch(ctx, func(tx *redis.Tx) error {
		for _, name := range names {
			uuid := fmt.Sprintf(x.playerInterimPrefix, uuid2.New().String())
			interStoreIDs = append(interStoreIDs, uuid)
			_, e := tx.ZInterStore(ctx, fmt.Sprintf(uuid), &redis.ZStore{
				Keys:      []string{pagination.OrderBy, fmt.Sprintf(x.playerByNamePrefix, name)},
				Aggregate: "SUM",
			}).Result()
			if e != nil {
				return fmt.Errorf("failed to create zinterstore: %w", e)
			}
		}
		resultID = fmt.Sprintf(x.playerInterimPrefix, uuid2.New().String())
		interStoreIDs = append(interStoreIDs, resultID)
		_, e := tx.ZInterStore(ctx, resultID, &redis.ZStore{
			Keys:      interStoreIDs[:len(interStoreIDs)-1],
			Aggregate: "SUM",
		}).Result()
		if e != nil {
			return fmt.Errorf("failed to create zunionstore: %w", e)
		}
		return nil
	}, pagination.OrderBy)
	if err != nil {
		return nil, fmt.Errorf("failed to search players: %w", err)
	}

	defer func(ids []string) {
		_, err := x.Del(ctx, ids...).Result()
		if err != nil {
			log.Printf("failed to clean up interim ids %v: %v", ids, err)
		}
	}(interStoreIDs)

	return x.searchByIndexAndPagination(ctx, resultID, pagination, withoutLimit)
}

func (x *playerRepository) searchByIndexAndPagination(
	ctx context.Context,
	index string,
	pagination model.PlayerPagination,
	withoutLimit bool,
) (*model.PlayersResponse, error) {
	var (
		ids     []string
		results []interface{}
		total   int64
	)

	err := x.Watch(ctx, func(tx *redis.Tx) error {
		var e error
		rangeFunc := tx.ZRevRange
		if pagination.ASC {
			rangeFunc = tx.ZRange
		}
		total, e = tx.ZCount(ctx, index, "-inf", "inf").Result()
		if e != nil {
			return fmt.Errorf("failed to zcount: %w", e)
		}
		start, stop := int64(pagination.Offset), int64(pagination.Offset+pagination.Limit-1)
		if withoutLimit {
			start, stop = 0, -1
		}
		ids, e = rangeFunc(ctx, index, start, stop).Result()
		if e != nil {
			return fmt.Errorf("failed to order by %s: %w", pagination.OrderBy, e)
		}
		if len(ids) == 0 {
			return ErrNotFound
		}
		results, e = tx.HMGet(ctx, x.playerRecords, ids...).Result()
		if e != nil {
			return fmt.Errorf("failed to hmget: %w", e)
		}
		return nil
	}, pagination.OrderBy, x.playerRecords)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to search players: %w", err)
	}

	players := make([]*model.Player, len(results))
	for key, result := range results {
		player := new(model.Player)
		data := []byte(result.(string))
		err = json.Unmarshal(data, player)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %w", err)
		}
		players[key] = player
	}

	return &model.PlayersResponse{
		Players: players,
		Total:   int(total),
		Offset:  pagination.Offset,
		Limit:   pagination.Limit,
	}, nil
}
