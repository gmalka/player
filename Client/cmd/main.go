package main

import (
	"log"

	"github.com/gmalka/Client/pkg/MusicPlayer"
	"github.com/gmalka/Client/pkg/songsManager"
	"github.com/gmalka/Client/pkg/fileManager"
	"github.com/gmalka/Client/internal/transport/grpc"
	"github.com/gmalka/Client/internal/service"
)

func main() {
	ch := make(chan byte, 1)
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
	uploadService, err := grpc.NewGrpcClient("localhost", "9879")
	if err != nil {
		log.Fatal(err)
	}

	manager := songsManager.NewSongManager(mp3FileManager, uploadService)

	controller := service.NewController(player, manager)
	controller.Run()
}