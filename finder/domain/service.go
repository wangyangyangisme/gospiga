package domain

import (
	"context"
)

type Service struct {
	db DB
}

func NewService(db DB) *Service {
	return &Service{db}
}

func (s *Service) IndexRecipe(ctx context.Context, recipe *Recipe) error {
	return nil
}
