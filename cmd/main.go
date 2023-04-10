package main

import (
	"context"
	"log"

	"github.com/gmalka/Player/pkg/Player"
	"github.com/gmalka/Player/pkg/controller"
	"github.com/gmalka/Player/pkg/fileManager"
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
	controller := controller.NewController(ctx, player, fm)
}