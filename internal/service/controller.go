package service

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"

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

type myController struct {
	ctx		context.Context
	Player
	MusicFileManager
}

func NewController(ctx context.Context, MusicPlayer Player, musicFileManager MusicFileManager) myController {
	return myController{ctx: ctx, Player: MusicPlayer, MusicFileManager: musicFileManager}
}

func (c myController) Run() {
	var (
		command []byte
	)

	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter command:")
		b, prefix, err := r.ReadLine()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for prefix == true {
			buf := []byte{}
			buf, prefix, _ = r.ReadLine()
			b = append(b, buf...)
		}
		n := bytes.IndexByte(b, ' ')
		if n == -1 {
			command = b
		} else {
			command = b[:n]
		}

		switch string(command) {
		case "Play", "play":
			if n != -1 {
				fmt.Println("incorrect command")
				continue
			}
			c.Player.Play()
		case "next", "nextsong", "Next song", "Next Song", "Nextsong":
			if n != -1 {
				fmt.Println("incorrect command")
				continue
			}
			c.Player.NextSong()
		case "pre", "presong", "Pre song", "Pre Song", "Presong":
			if n != -1 {
				fmt.Println("incorrect command")
				continue
			}
			c.Player.PreSong()
		case "pause", "Pause", "stop", "Stop":
			if n != -1 {
				fmt.Println("incorrect command")
				continue
			}
			c.Player.Pause()
		}
	}
}
