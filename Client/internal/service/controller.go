package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)


type myController struct {
	player player
	songmanager songmanager
}

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
	Delete(id int) error
	DeleteLocal(name string) error
	SaveLocal(name string) error
	GetAllLocal() ([]string, error)
	GetAllRemote() ([]string, error)
}

func NewController(MusicPlayer player, musicFileManager songmanager) myController {
	return myController{player: MusicPlayer, songmanager: musicFileManager}
}

func (c *myController) SetVolume(str string) error {
	v, err := strconv.Atoi(str)
	if err != nil {
		return errors.New(fmt.Sprintf("Incorrect input value: %s", str))
	}
	err = c.player.SetVolume(v)
	if err != nil {
		log.Println(err)
	}
	return nil
}
func (c *myController) GetCurrent() string {
	return c.songmanager.GetCurrent()
}

func (c *myController) AddSong(str string) error {
	err := c.songmanager.Add(str)
	if err != nil {
		return errors.New(fmt.Sprintf("Cant find song %s", str))
	}
	return nil
}

func (c *myController) GetPlayList() ([]string, error) {
	songs := c.songmanager.GetPlayList()
	if songs == nil {
		return nil, errors.New("Play list is empty")
	}
	return songs, nil
}

func (c *myController) DeleteLocal(str string) error {
	if str != "" {
		err := c.songmanager.DeleteLocal(str)
		if err != nil {
			return err
		}
	} else {
		err := c.songmanager.DeleteLocal("")
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *myController) DeleteSong(id int) error {
	err := c.songmanager.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *myController) GetAllSongs(str string) ([][]string, error) {
	if str == "" {
		local, err := c.songmanager.GetAllLocal()
		if err != nil {
			return nil, err
		}
		remote, err := c.songmanager.GetAllRemote()
		if err != nil {
			return nil, err
		} else {
			return [][]string{local, remote}, nil
		}
	} else if str == "local" || str == "Local" {
		local, err := c.songmanager.GetAllLocal()
		if err != nil {
			return nil, err
		}
		return [][]string{local, nil}, nil
	} else if str == "remote" || str == "Remote" {
		remote, err := c.songmanager.GetAllRemote()
		if err != nil {
			return nil, err
		}
		return [][]string{nil, remote}, nil
	}
	return nil, errors.New("Incorrect command")
}

func (c *myController) SaveSong(str string) error {
	if str != "" {
		err := c.songmanager.SaveLocal(str)
		if err != nil {
			return err
		}
	} else {
		err := c.songmanager.SaveLocal("")
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *myController) StopSong() {
	c.player.Stop()
}

func (c *myController) PauseSong() {
	c.player.Pause()
}

func (c *myController) PreSong() error {
	data, err := c.songmanager.Pre()
	if err != nil {
		return err
	} else {
		c.player.Load(data)
	}
	return nil
}

func (c *myController) NextSong() error {
	data, err := c.songmanager.Next()
	if err != nil {
		return err
	} else {
		c.player.Load(data)
	}
	return nil
}

func (c *myController) PlaySong(str string) error {
	if c.player.IsPlaying() {
		c.player.Play()
	} else {
		data, err := c.songmanager.Get(str)
		if err != nil {
			return err
		} else {
			c.player.Load(data)
			c.player.Play()
		}
	}
	return nil
}