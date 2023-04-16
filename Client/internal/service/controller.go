package service

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
)

type player interface {
	Play()
	Pause()
	Stop()
	SetVolume(v int) error
 	Load(data []byte) error
	IsPlaying() bool
	GetSongInfo() string
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
	GetAllLocal() ([]string, error)
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
func (c *myController) Run() {
	var (
		command []byte
	)

	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command:     ")
		b, prefix, err := r.ReadLine()
		if err != nil {
			log.Println(err)
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
		case "add", "addsong", "Add", "AddSong", "addSong":
			c.AddSong(b, n)
		case "Play", "play":
			c.PlaySong(b, n)
		case "next", "Next", "nextsong", "Nextsong", "NextSong":
			c.NextSong(b, n)
		case "pre", "Pre", "presong", "Presong", "PreSong":
			c.PreSong(b, n)
		case "pause", "Pause":
			c.PauseSong(b, n)
		case "stop", "Stop":
			c.StopSong(b, n)
		case "Save", "SaveLocal", "save", "savelcoal", "saveLocal":
			c.SaveSong(b, n)
		case "getall", "Getall", "getAll", "GetAll":
			c.GetAllSongs(b, n)
		case "delete", "Delete":
			c.DeleteSong(b, n)
		case "deletelocal", "Deletelocal", "DeleteLocal", "deleteLocal":
			c.DeleteLocal(b, n)
		case "playlist", "list", "List", "Playlist", "PlayList", "playList":
			c.GetPlayList(b, n)
		case "volume", "Volume", "vol", "Vol":
			c.SetVolume(b, n)
		case "status", "Status", "info":
			fmt.Printf("#>Name:%s | Playing: %t\n#>Time: %s\n", c.songmanager.GetCurrent(), c.player.IsPlaying(),  c.player.GetSongInfo())
		case "exit", "Exit":
			c.player.Stop()
			return
		default:
			fmt.Printf("#?>Unknown command: \"%s\"\n", string(b))
		}
	}
}

func (c *myController) SetVolume(b []byte, n int) {
	if n == -1 {
		fmt.Printf("#?>Incorrect command\n")
		return
	}

	v, err := strconv.Atoi(string(b[n+1:]))
	if err != nil {
		fmt.Printf("#?>Incorrect command \"%s\"\n", string(b[n+1:]))
		return
	}
	err = c.player.SetVolume(v)
	if err != nil {
		log.Println(err)
	}
}

func (c *myController) AddSong(b []byte, n int) {
	if n == -1 {
		fmt.Printf("#?>Incorrect command\n")
		return
	}
	
	err := c.songmanager.Add(string(b[n+1:]))
	if err != nil {
		log.Println(err)
	}
}

func (c *myController) GetPlayList(b []byte, n int) {
	if n != -1 {
		fmt.Printf("#?>Incorrect command\n")
		return
	}
	songs := c.songmanager.GetPlayList()
	if songs == nil {
		log.Println("Play list is empty")
	} else {
		var (
			i int
			s string
		)
		for i, s = range songs {
			if (i + 1) % 3 != 0 {
				fmt.Printf("%3d: %-20s | ", i + 1, s)
			} else {
				fmt.Printf("%3d: %-20s\n", i + 1, s)
			} 
		}
		if (i + 1) % 3 != 0 {
			fmt.Println()
		}
	}
}

func (c *myController) DeleteLocal(b []byte, n int) {
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
}

func (c *myController) DeleteSong(b []byte, n int) {
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
}

func (c *myController) GetAllSongs(b []byte, n int) {
	if n == -1 {
		local, err := c.songmanager.GetAllLocal()
		if err != nil {
			log.Println(err)
			return
		}
		remote, err := c.songmanager.GetAllRemote()
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println("#>Songs on Server: ")
			var (
				i int
				s string
			)
			for i, s = range remote {
				if (i + 1) % 3 != 0 {
					fmt.Printf("%3d: %-20s | ", i + 1, s)
				} else {
					fmt.Printf("%3d: %-20s\n", i + 1, s)
				}
			}
			if (i + 1) % 3 != 0 {
				fmt.Println()
			}
			fmt.Println("#>Local songs: ")
			for i, s = range local {
				if (i + 1) % 3 != 0 {
					fmt.Printf("%3d: %-20s | ", i + 1, s)
				} else {
					fmt.Printf("%3d: %-20s\n", i + 1, s)
				}
			}
			if (i + 1) % 3 != 0 {
				fmt.Println()
			}
		}
	} else if comm2 := string(b[n+1:]); comm2 == "local" || comm2 == "Local" {
		local, err := c.songmanager.GetAllLocal()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("#>Local songs: ")
		var (
			i int
			s string
		)
		for i, s = range local {
			if (i + 1) % 3 != 0 {
				fmt.Printf("%3d: %-20s | ", i + 1, s)
			} else {
				fmt.Printf("%3d: %-20s\n", i + 1, s)
			}
		}
		if (i + 1) % 3 != 0 {
			fmt.Println()
		}
	} else if comm2 == "remote" || comm2 == "Remote" {
		remote, err := c.songmanager.GetAllRemote()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("#>Songs on Server: ")
		var (
			i int
			s string
		)
		for i, s = range remote {
			if (i + 1) % 3 != 0 {
				fmt.Printf("%3d: %-20s | ", i + 1, s)
			} else {
				fmt.Printf("%3d: %-20s\n", i + 1, s)
			}
		}
		if (i + 1) % 3 != 0 {
			fmt.Println()
		}
	}
}

func (c *myController) SaveSong(b []byte, n int) {
	if n == -1 {
		err := c.songmanager.SaveLocal("")
		if err != nil {
			log.Println(err)
		}
	} else {
		err := c.songmanager.SaveLocal(string(b[n+1:]))
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *myController) StopSong(b []byte, n int) {
	if n != -1 {
		fmt.Printf("#?>Incorrect command\n")
		return
	}
	c.player.Stop()
}

func (c *myController) PauseSong(b []byte, n int) {
	if n != -1 {
		fmt.Printf("#?>Incorrect command\n")
		return
	}
	c.player.Pause()
}

func (c *myController) PreSong(b []byte, n int) {
	if n != -1 {
		fmt.Printf("#?>Incorrect command\n")
		return
	}
	data, err := c.songmanager.Pre()
	if err != nil {
		log.Println(err)
	} else {
		c.player.Load(data)
	}
}

func (c *myController) NextSong(b []byte, n int) {
	if n != -1 {
		fmt.Printf("#?>Incorrect command\n")
		return
	}
	data, err := c.songmanager.Next()
	if err != nil {
		log.Println(err)
	} else {
		c.player.Load(data)
	}
}

func (c *myController) PlaySong(b []byte, n int) {
	if n != -1 {
		fmt.Printf("#?>Incorrect command\n")
		return
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
}