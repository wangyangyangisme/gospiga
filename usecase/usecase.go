package usecase

import (
	"context"
	"fmt"

	gostreamer "github.com/kind84/gospiga/streamer"
)

type App struct {
	service  Service
	db       DB
	streamer Streamer
	provider Provider
}

func NewApp(ctx context.Context, service Service, db DB, streamer Streamer, provider Provider) *App {
	app := &App{
		service:  service,
		db:       db,
		streamer: streamer,
		provider: provider,
	}

	// start streamer to listen for new recipes.
	go app.readNewRecipes(ctx)

	return app
}

func (a *App) readNewRecipes(ctx context.Context) {
	msgChan := make(chan gostreamer.Message)
	exitChan := make(chan struct{})

	args := &gostreamer.StreamArgs{
		Stream:   "new-recipes",
		Group:    "usecase",
		Consumer: "usecase",
	}
	a.streamer.ReadGroup(ctx, args, msgChan, exitChan)

	for {
		select {
		case msg := <-msgChan:
			recipeID, ok := msg.Payload.(string)
			if !ok {
				fmt.Println("cannot read recipe ID from message.")
			}
			fmt.Printf("Got message for a new recipe ID %s\n", recipeID)
			// call datocms to get the full recipe
			r, err := a.provider.GetRecipe(ctx, recipeID)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// save message
			fmt.Println(*r)
			a.service.SaveRecipe(ctx, r)
			// ack streamer

		case <-ctx.Done():
			// time to exit
		}
	}
}
