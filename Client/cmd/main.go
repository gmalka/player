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
	ch := make(chan byte)
	defer close(ch)

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

	player, err := MusicPlayer.NewMp3Player(ch, manager)
	if err != nil {
		log.Fatal(err)
	}

	controller := service.NewController(player, manager)
	controller.Run()
}