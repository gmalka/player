package main

import (
	"io"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	cli "github.com/gmalka/Client/internal/CLI"
	"github.com/gmalka/Client/internal/service"
	"github.com/gmalka/Client/internal/transport/grpc"
	"github.com/gmalka/Client/pkg/MusicPlayer"
	"github.com/gmalka/Client/pkg/fileManager"
	"github.com/gmalka/Client/pkg/songsManager"
	"github.com/spf13/viper"
)


type mod interface {
	Height() int
	Spacing() int
	Update(msg tea.Msg, m *list.Model) tea.Cmd
	Render(w io.Writer, m list.Model, index int, listItem list.Item)
}


func main() {
	ch := make(chan byte)
	defer close(ch)

	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	mp3FileManager, err := fileManager.NewMusicFileManager(fileManager.DefaultPath)
	if err != nil {
		log.Fatal(err)
	}

	uploadService, err := grpc.NewGrpcClient(viper.GetString("ip"), viper.GetString("port"))
	if err != nil {
		log.Fatal(err)
	}
	
	manager := songsManager.NewSongManager(mp3FileManager, uploadService)

	player, err := MusicPlayer.NewMp3Player(ch, manager)
	if err != nil {
		log.Fatal(err)
	}

	controller := service.NewController(player, manager)
	
	cli.RunModel(&controller, []string{
		"Add", "Play", "Pause", "Set Volume", "Next", "Pre", "Playlist", "Get all songs", "Delete", "Stop", "Delete from local storage", "Save song in local storage"})
}