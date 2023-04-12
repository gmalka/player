package service

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

type player interface {
	Play()
	Pause()
	Stop()
	Load(data []byte) error
	IsPlaying() bool
}

type songmanager interface {
	Get(name string) ([]byte, error)
	Add(name string) error
	Next() ([]byte, error)
	Pre() ([]byte, error)
	GetCurrent() string
	GetPlayList() []string
	Delete(name string) error
	DeleteLocal(name string) error
	SaveLocal(name string) error
	GetAllLocal() []string
	GetAllRemote() ([]string, error)
}

type myController struct {
	player player
	songmanager songmanager
}

func NewController(MusicPlayer player, musicFileManager songmanager) myController {
	return myController{player: MusicPlayer, songmanager: musicFileManager}
}

//TODO: раскидать тела case по функциям
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
			if c.player.IsPlaying() {
				c.player.Play()
			} else {
				data, err := c.songmanager.Get("")
				if err != nil {
					log.Println(err)
				} else {
					c.player.Load(data)
					c.player.Play()
				}
			}
			c.player.Play()
		case "next", "Next", "nextsong", "Nextsong", "NextSong":
			if n != -1 {
				fmt.Println("incorrect command")
				continue
			}
			data, err := c.songmanager.Pre()
			if err != nil {
				log.Println(err)
			} else {
				c.player.Load(data)
				c.player.Play()
			}
		case "pre", "Pre", "presong", "Presong", "PreSong":
			if n != -1 {
				fmt.Println("incorrect command")
				continue
			}
			data, err := c.songmanager.Next()
			if err != nil {
				log.Println(err)
			} else {
				c.player.Load(data)
				c.player.Play()
			}
		case "pause", "Pause", "stop", "Stop":
			if n != -1 {
				fmt.Println("incorrect command")
				continue
			}
			c.player.Pause()
		case "Save", "SaveLocal", "save", "savelcoal", "saveLocal":
			if n == -1 {
				err = c.songmanager.SaveLocal("")
				if err != nil {
					log.Println(err)
				}
			} else {
				err = c.songmanager.SaveLocal(string(b[n+1:]))
				if err != nil {
					log.Println(err)
				}
			}
		case "getall", "Getall", "getAll", "GetAll":
			if n == -1 {
				local := c.songmanager.GetAllLocal()
				remote, err := c.songmanager.GetAllRemote()
				if err != nil {
					log.Println(err)
				} else {
					fmt.Println("Songs on Server: ")
					for i, s := range remote {
						if i % 3 != 0 {
							fmt.Printf("%100d: %-30s | ", i, s)
						} else {
							fmt.Printf("%100d: %-30s\n", i, s)
						}
					}
					fmt.Println("Local songs: ")
					for i, s := range local {
						if i % 3 != 0 {
							fmt.Printf("%100d: %-30s | ", i, s)
						} else {
							fmt.Printf("%100d: %-30s\n", i, s)
						}
					}
				}
			} else if comm2 := string(b[n+1:]); comm2 == "local" || comm2 == "Local" {
				local := c.songmanager.GetAllLocal()
				fmt.Println("Local songs: ")
				for i, s := range local {
					if i % 3 != 0 {
						fmt.Printf("%100d: %-30s | ", i, s)
					} else {
						fmt.Printf("%100d: %-30s\n", i, s)
					}
				}
			} else if comm2 == "remote" || comm2 == "Remote" {
				remote, err := c.songmanager.GetAllRemote()
				if err != nil {
					log.Println(err)
				}
				fmt.Println("Songs on Server: ")
				for i, s := range remote {
					if i % 3 != 0 {
						fmt.Printf("%100d: %-30s | ", i, s)
					} else {
						fmt.Printf("%100d: %-30s\n", i, s)
					}
				}
			}
		case "delete", "Delete":
			if n != -1 {
				err := c.songmanager.Delete(string(b[n+1:]))
				if err != nil {
					log.Println(err)
				}
			} else {
				err := c.songmanager.Delete("")
				if err != nil {
					log.Println(err)
				}
			}
		case "deletelocal", "Deletelocal", "DeleteLocal", "deleteLocal":
			if n != -1 {
				err := c.songmanager.DeleteLocal(string(b[n+1:]))
				if err != nil {
					log.Println(err)
				}
			} else {
				err := c.songmanager.DeleteLocal("")
				if err != nil {
					log.Println(err)
				}
			}
		case "playlist", "list", "List", "Playlist", "PlayList", "playList":
			songs := c.songmanager.GetPlayList()
			if songs == nil {
				fmt.Println("Play list is empty")
			} else {
				for i, s := range songs {
					if i % 3 != 0 {
						fmt.Printf("%100d: %-30s | ", i, s)
					} else {
						fmt.Printf("%100d: %-30s\n", i, s)
					} 
				}
			}
		case "status", "Status":
			fmt.Sprintf("Name:%s Status: %b", c.songmanager.GetCurrent(), c.player.IsPlaying())
		}
	}
}
