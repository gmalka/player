package controller

import (
	"context"

	"github.com/hajimehoshi/go-mp3"
)

type Player interface {
	Play()
	Pause()
	NextSong()
	PreSong()
	Add(decorded *mp3.Decoder)
}

type MusicFileManager interface {
	Add(name string, input []byte) error
	Open(name string) ([]byte, error)
	GetAll() []string
	Delete(name string) error
}

type Controller interface {
	ManageInputCommands()
}

type myController struct {
	ctx		context.Context
	Player
	MusicFileManager
}

func NewController(ctx context.Context, MusicPlayer Player, musicFileManager MusicFileManager) myController {
	return myController{ctx: ctx, Player: MusicPlayer, MusicFileManager: musicFileManager}
}
