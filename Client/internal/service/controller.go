package service

import (
	"errors"
	"fmt"
	"strconv"
)

type myController struct {
	preSong     string
	player      player
	songmanager songmanager
}

type Controller interface {
	SetVolume(str string) error
	AddSong(str string) error
	GetPlayList() ([]string, error)
	DeleteLocal(str string) error
	DeleteSong(id int) error
	GetAllSongs(str string) ([][]string, error)
	SaveSong(str string) error
	StopSong()
	PauseSong()
	PreSong() error
	NextSong() error
	PlaySong(str string) error
	GetCurrent() string
}

type player interface {
	Play()
	Pause()
	Stop()
	SetVolume(v int) error
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
	Delete(id int) error
	DeleteLocal(name string) error
	SaveLocal(name string) error
	GetAllLocal() ([]string, error)
	GetAllRemote() ([]string, error)
}

func NewController(MusicPlayer player, musicFileManager songmanager) Controller {
	return &myController{player: MusicPlayer, songmanager: musicFileManager}
}

func (c *myController) SetVolume(str string) error {
	v, err := strconv.Atoi(str)
	if err != nil {
		return errors.New(fmt.Sprintf("Incorrect input value: %s", str))
	}
	err = c.player.SetVolume(v)
	if err != nil {
		return err
	}
	return nil
}

func (c *myController) GetCurrent() string {
	if c.player.IsPlaying() {
		return c.preSong
	}
	return c.songmanager.GetCurrent()
}

func (c *myController) AddSong(str string) error {
	err := c.songmanager.Add(str)

	if err != nil {
		return errors.New(fmt.Sprintf("Cant find song %s", str))
	}
	if c.preSong == "" {
		c.preSong = str
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
	err := c.songmanager.SaveLocal(str)
	if err != nil {
		return err
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
		c.preSong = c.songmanager.GetCurrent()
		return c.player.Load(data)
	}
}

func (c *myController) NextSong() error {
	data, err := c.songmanager.Next()
	if err != nil {
		return err
	} else {
		c.preSong = c.songmanager.GetCurrent()
		return c.player.Load(data)
	}
}

func (c *myController) PlaySong(str string) error {
	if c.player.IsPlaying() {
		c.player.Play()
	} else {
		data, err := c.songmanager.Get(str)
		if err != nil {
			return err
		} else {
			err := c.player.Load(data)
			if err != nil {
				return err
			}
			c.player.Play()
		}
	}
	return nil
}
