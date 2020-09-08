package graph

import (
	"the-rush/repository"
)

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	repo *repository.Repository
}

func NewResolver(repository *repository.Repository) *Resolver {
	return &Resolver{
		repo: repository,
	}
}
