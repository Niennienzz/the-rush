package model

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"the-rush/constant"
)

type Player struct {
	ID                            string       `json:"id"`
	CreatedAt                     time.Time    `json:"createdAt"`
	Name                          string       `json:"name"`
	Team                          string       `json:"team"`
	Position                      Position     `json:"position"`
	RushingAttempts               int          `json:"rushingAttempts"`
	RushingAttemptsPerGameAverage float64      `json:"rushingAttemptsPerGameAverage"`
	TotalRushingYards             int          `json:"totalRushingYards"`
	RushingAverageYardsPerAttempt float64      `json:"rushingAverageYardsPerAttempt"`
	RushingYardsPerGame           float64      `json:"rushingYardsPerGame"`
	TotalRushingTouchdowns        int          `json:"totalRushingTouchdowns"`
	LongestRush                   *LongestRush `json:"longestRush"`
	RushingFirstDowns             int          `json:"rushingFirstDowns"`
	RushingFirstDownsPercentage   float64      `json:"rushingFirstDownsPercentage"`
	Rushing20PlusYardsEach        int          `json:"rushing20PlusYardsEach"`
	Rushing40PlusYardsEach        int          `json:"rushing40PlusYardsEach"`
	RushingFumbles                int          `json:"rushingFumbles"`
}

func (x Player) marshalCSV() []byte {
	return []byte(fmt.Sprintf("%s,%d,%d,%d,%t,%s,%s,%d,%v,%v,%d,%v,%d,%d,%d\n",
		x.Name, x.TotalRushingYards, x.TotalRushingTouchdowns, x.LongestRush.Value, x.LongestRush.IsTouchdown,
		x.Team, x.Position.String(), x.RushingAttempts, x.RushingAttemptsPerGameAverage,
		x.RushingYardsPerGame, x.RushingFirstDowns, x.RushingFirstDownsPercentage,
		x.Rushing20PlusYardsEach, x.Rushing40PlusYardsEach, x.RushingFumbles,
	))
}

type PlayersCSV []*Player

func (x PlayersCSV) Marshal() []byte {
	buf := new(bytes.Buffer)
	header := "Player Name,Total Rushing Yards,Total Rushing Touchdowns,Longest Rush,Longest Rush Is Touchdown,Team,Position,Rushing Attempts,Rushing Attempts Per Game Average,Rushing Yards Per Game,Rushing First Downs,Rushing First Downs Percentage,Rushing 20+ Yards Each,Rushing 40+ Yards Each,Rushing Fumbles\n"
	buf.Write([]byte(header))
	for _, player := range x {
		buf.Write(player.marshalCSV())
	}
	return buf.Bytes()
}

type PlayerInput struct {
	Name                          string      `json:"Player"`
	Team                          string      `json:"Team"`
	Position                      Position    `json:"Pos"`
	RushingAttempts               int         `json:"Att"`
	RushingAttemptsPerGameAverage float64     `json:"Att/G"`
	TotalRushingYards             interface{} `json:"Yds"`
	RushingAverageYardsPerAttempt float64     `json:"Avg"`
	RushingYardsPerGame           float64     `json:"Yds/G"`
	TotalRushingTouchdowns        int         `json:"TD"`
	LongestRush                   interface{} `json:"Lng"`
	RushingFirstDowns             int         `json:"1st"`
	RushingFirstDownsPercentage   float64     `json:"1st%"`
	Rushing20PlusYardsEach        int         `json:"20+"`
	Rushing40PlusYardsEach        int         `json:"40+"`
	RushingFumbles                int         `json:"FUM"`
}

func (x PlayerInput) ToPlayer() (*Player, error) {
	var (
		yds int
		err error
	)
	switch v := x.TotalRushingYards.(type) {
	case float64:
		yds = int(v)
	case string:
		yds, err = strconv.Atoi(strings.ReplaceAll(v, ",", ""))
		if err != nil {
			return nil, fmt.Errorf("failed to parse totalRushingYards: %w", err)
		}
	}

	// A 'T' represents a touchdown occurred.
	longestRush := new(LongestRush)
	switch v := x.LongestRush.(type) {
	case float64:
		longestRush.Value = int(v)
	case string:
		if strings.HasSuffix(v, "T") {
			longestRush.IsTouchdown = true
			v = v[:len(v)-1]
		}
		value, err := strconv.Atoi(strings.ReplaceAll(v, ",", ""))
		if err != nil {
			return nil, fmt.Errorf("failed to parse longestRushValue: %w", err)
		}
		longestRush.Value = value
	}

	return &Player{
		ID:                            uuid.New().String(),
		CreatedAt:                     time.Now(),
		Name:                          x.Name,
		Team:                          x.Team,
		Position:                      x.Position,
		RushingAttempts:               x.RushingAttempts,
		RushingAttemptsPerGameAverage: x.RushingAttemptsPerGameAverage,
		TotalRushingYards:             yds,
		RushingAverageYardsPerAttempt: x.RushingAverageYardsPerAttempt,
		RushingYardsPerGame:           x.RushingYardsPerGame,
		TotalRushingTouchdowns:        x.TotalRushingTouchdowns,
		LongestRush:                   longestRush,
		RushingFirstDowns:             x.RushingFirstDowns,
		RushingFirstDownsPercentage:   x.RushingFirstDownsPercentage,
		Rushing20PlusYardsEach:        x.Rushing20PlusYardsEach,
		Rushing40PlusYardsEach:        x.Rushing40PlusYardsEach,
		RushingFumbles:                x.RushingFumbles,
	}, nil
}

type PlayerPagination struct {
	Offset  int
	Limit   int
	OrderBy string
	ASC     bool
}

func PlayerPaginationFromArgs(args PlayersArgs) (*PlayerPagination, error) {
	var (
		offset   = 0
		limit    = 25
		orderBy  = constant.RedisIndexKeyPlayerByCreatedAt.String()
		orderAsc = true
	)

	if args.Order != nil {
		switch args.Order.OrderBy {
		case PlayersArgsOrderByCreatedAt:
			orderBy = constant.RedisIndexKeyPlayerByCreatedAt.String()
		case PlayersArgsOrderByLongestRush:
			orderBy = constant.RedisIndexKeyPlayerByLongestRush.String()
		case PlayersArgsOrderByTotalRushingTouchdowns:
			orderBy = constant.RedisIndexKeyPlayerByTotalRushingTouchdowns.String()
		case PlayersArgsOrderByTotalRushingYards:
			orderBy = constant.RedisIndexKeyPlayerByTotalRushingYards.String()
		default:
			return nil, errors.New("invalid order by in players args")
		}

		switch args.Order.Order {
		case OrderAsc:
			orderAsc = true
		case OrderDesc:
			orderAsc = false
		default:
			return nil, errors.New("invalid order in players args")
		}
	}

	if args.Page != nil {
		offset = args.Page.Offset
		if pageLimit := args.Page.Limit; pageLimit < constant.MaxPageLimit {
			limit = pageLimit
		} else {
			limit = constant.MaxPageLimit
		}
	}

	return &PlayerPagination{
		Offset:  offset,
		Limit:   limit,
		OrderBy: orderBy,
		ASC:     orderAsc,
	}, nil
}
