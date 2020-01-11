package domain

import (
	"context"
)

type DB interface {
	SaveRecipe(context.Context, *Recipe) error
	GetRecipeByID(context.Context, string) (*Recipe, error)
	IDSaved(context.Context, string) (bool, error)
}
