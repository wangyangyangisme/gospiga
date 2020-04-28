// +build integration

package dgraph

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"

	"github.com/kind84/gospiga/pkg/errors"
	"github.com/kind84/gospiga/server/domain"
)

var db *DB

func init() {
	d, err := grpc.Dial("alpha:9080", grpc.WithInsecure())
	if err != nil {
		panic(fmt.Errorf("failed to connect to dgraph: %w", err))
	}

	dgraph := dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
	db = &DB{dgraph}
}

func TestSaveRecipe(t *testing.T) {
	recipe := getTestRecipe()

	tests := []struct {
		name        string
		setup       func(ctx context.Context, db *DB) error
		recipe      *domain.Recipe
		expectedErr error
	}{
		{
			name:   "save new recipe",
			recipe: recipe,
		},
		{
			name: "don't save same xid again",
			setup: func(ctx context.Context, db *DB) error {
				return db.SaveRecipe(ctx, recipe)
			},
			recipe:      recipe,
			expectedErr: errors.ErrDuplicateID{ID: recipe.ID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			require := require.New(t)
			if tt.setup != nil {
				err := tt.setup(ctx, db)
				require.NoError(err)
			}

			err := db.SaveRecipe(ctx, tt.recipe)

			if tt.expectedErr != nil {
				require.Error(err)
				require.IsType(errors.ErrDuplicateID{ID: recipe.ExternalID}, err)
			} else {
				require.NoError(err)
			}
			r, err := db.GetRecipeByID(ctx, tt.recipe.ExternalID)
			require.NoError(err)
			require.NotNil(r)
			assert.Equal(t, r.ExternalID, tt.recipe.ExternalID)
			assert.Equal(t, r.Title, tt.recipe.Title)
			n, err := db.CountRecipes(ctx)
			require.NoError(err)
			assert.Equal(t, n, 1)
			err = db.DeleteRecipe(ctx, tt.recipe.ExternalID)
			require.NoError(err)
		})
	}
}

func TestUpdateRecipe(t *testing.T) {
	recipe := getTestRecipe()
	recipe2 := getTestRecipe()
	recipe3 := getTestRecipe()
	recipe2.Title = "update"
	recipe3.Ingredients[0].Quantity = 10

	tests := []struct {
		name    string
		recipe  *domain.Recipe
		setup   func(ctx context.Context, db *DB) error
		assert  func(ctx context.Context, db *DB, t *testing.T)
		cleanup func(ctx context.Context, db *DB) error
	}{
		{
			name:   "recipe not found not updated",
			recipe: recipe2,
			assert: func(ctx context.Context, db *DB, t *testing.T) {
				require := require.New(t)
				assert := assert.New(t)
				n, err := db.CountRecipes(ctx)
				require.NoError(err)
				assert.Equal(0, n)
			},
		},
		{
			name:   "recipe found scalar field updated",
			recipe: recipe2,
			setup: func(ctx context.Context, db *DB) error {
				return db.SaveRecipe(ctx, recipe)
			},
			assert: func(ctx context.Context, db *DB, t *testing.T) {
				require := require.New(t)
				assert := assert.New(t)
				r, err := db.GetRecipeByID(ctx, recipe2.ExternalID)
				require.NoError(err)
				require.NotNil(r)
				n, err := db.CountRecipes(ctx)
				require.NoError(err)
				assert.Equal(1, n)
				assert.Equal(r.ExternalID, recipe2.ExternalID)
				assert.Equal(r.Title, recipe2.Title)
			},
			cleanup: func(ctx context.Context, db *DB) error {
				return db.DeleteRecipe(ctx, recipe2.ExternalID)
			},
		},
		{
			name:   "recipe found ingredient found field updated",
			recipe: recipe3,
			setup: func(ctx context.Context, db *DB) error {
				return db.SaveRecipe(ctx, recipe)
			},
			assert: func(ctx context.Context, db *DB, t *testing.T) {
				require := require.New(t)
				assert := assert.New(t)
				r, err := db.GetRecipeByID(ctx, recipe3.ExternalID)
				require.NoError(err)
				require.NotNil(r)
				n, err := db.CountRecipes(ctx)
				require.NoError(err)
				assert.Equal(1, n)
				assert.Equal(recipe3.ExternalID, r.ExternalID)
				if qs, ok := r.Ingredients[0].Quantity.(string); assert.True(ok) {
					q, err := strconv.Atoi(qs)
					assert.NoError(err)
					assert.Equal(recipe3.Ingredients[0].Quantity, q)
				}
				assert.Equal(recipe3.Ingredients[0].Name, r.Ingredients[0].Name)
				assert.Equal(recipe3.Ingredients[0].UnitOfMeasure, r.Ingredients[0].UnitOfMeasure)
				assert.Equal(1, len(r.Ingredients))
			},
			cleanup: func(ctx context.Context, db *DB) error {
				return db.DeleteRecipe(ctx, recipe3.ExternalID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			require := require.New(t)
			if tt.setup != nil {
				err := tt.setup(ctx, db)
				require.NoError(err)
			}

			err := db.UpdateRecipe(ctx, tt.recipe)

			require.NoError(err)
			if tt.assert != nil {
				tt.assert(ctx, db, t)
			}
			if tt.cleanup != nil {
				err := tt.cleanup(ctx, db)
				require.NoError(err)
			}
		})
	}
}

func TestDeleteRecipe(t *testing.T) {
	recipe := getTestRecipe()

	err := db.SaveRecipe(context.Background(), recipe)
	require.NoError(t, err)

	err = db.DeleteRecipe(context.Background(), recipe.ExternalID)

	require.NoError(t, err)
	recipe, err = db.GetRecipeByID(context.Background(), recipe.ExternalID)
	require.NoError(t, err)
	require.Nil(t, recipe)
}

func TestGetRecipeByID(t *testing.T) {
	recipe := getTestRecipe()
	err := db.SaveRecipe(context.Background(), recipe)
	require.NoError(t, err)

	r, err := db.GetRecipeByID(context.Background(), "externalID")

	require.NoError(t, err)
	assert.Equal(t, recipe.Title, r.Title)
	assert.Equal(t, recipe.Subtitle, r.Subtitle)
	if assert.Equal(t, len(recipe.Ingredients), len(r.Ingredients)) {
		for i := 0; i < len(recipe.Ingredients); i++ {
			assert.Equal(t, recipe.Ingredients[i].Name, r.Ingredients[i].Name)
			assert.Equal(t, recipe.Ingredients[i].UnitOfMeasure, r.Ingredients[i].UnitOfMeasure)
			if qs, ok := r.Ingredients[i].Quantity.(string); assert.True(t, ok) {
				q, err := strconv.Atoi(qs)
				assert.NoError(t, err)
				assert.Equal(t, recipe.Ingredients[i].Quantity, q)
			}
		}
	}
	if assert.Equal(t, len(recipe.Steps), len(r.Steps)) {
		for i := 0; i < len(recipe.Steps); i++ {
			assert.Equal(t, recipe.Steps[i].Heading, r.Steps[i].Heading)
			assert.Equal(t, recipe.Steps[i].Body, r.Steps[i].Body)
		}
	}
	err = db.DeleteRecipe(context.Background(), recipe.ExternalID)
	require.NoError(t, err)
}

func getTestRecipe() *domain.Recipe {
	return &domain.Recipe{
		ExternalID:  "externalID",
		Title:       "title",
		Subtitle:    "subtitle",
		Description: "description",
		Conclusion:  "conclusion",
		MainImage:   &domain.Image{URL: "url"},
		Difficulty:  domain.DifficultyEasy,
		Cost:        domain.CostLow,
		Servings:    1,
		PrepTime:    1,
		CookTime:    1,
		Ingredients: []*domain.Ingredient{
			{
				Name:          "ingredient",
				Quantity:      1,
				UnitOfMeasure: "unitOfMeasure",
			},
		},
		Steps: []*domain.Step{
			{
				Heading: "heading",
				Body:    "body",
				Image:   &domain.Image{URL: "url"},
			},
		},
		Tags: []*domain.Tag{{TagName: "tagName"}},
		Slug: "test-recipe",
	}
}
