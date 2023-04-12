package main

import (
	"context"
	"log"

	"github.com/gmalka/Player/internal/app"
	"github.com/gmalka/Player/pkg/fileManager"
	"github.com/gmalka/Player/pkg/player"
)

func main() {
	ch := make(chan byte)
	player, err := Player.NewMp3Player(ch)
	defer close(ch)
	if err != nil {
		log.Fatal(err)
	}
	fm, err := fileManager.NewMusicFileManager(fileManager.DefaultPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	controller := app.NewController(ctx, player, fm)
}