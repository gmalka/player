package main

import (
	"context"
	"log"

	"github.com/gmalka/Player/internal/app"
	"github.com/gmalka/Player/internal/service"
	"github.com/gmalka/Player/internal/transport/grpc"
	"github.com/gmalka/Player/pkg/MusicPlayer"
	"github.com/gmalka/Player/pkg/fileManager"
	"github.com/gmalka/Player/pkg/songsManager"
)

func main() {
	ch := make(chan byte)
	defer close(ch)

	player, err := MusicPlayer.NewMp3Player(ch)
	if err != nil {
		log.Fatal(err)
	}

	mp3FileManager, err := fileManager.NewMusicFileManager(fileManager.DefaultPath)
	if err != nil {
		log.Fatal(err)
	}


	//TODO: создать файл конфигов
	uploadService, err := grpc.NewGrpcClient("localhost", "5762")
	if err != nil {
		log.Fatal(err)
	}

	manager := songsManager.NewSongManager(mp3FileManager, uploadService)

	controller := service.NewController(player, manager)
	controller.Run()
}