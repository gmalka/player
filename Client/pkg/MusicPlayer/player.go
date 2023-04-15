package MusicPlayer

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

type Player interface {
	Play()
	Pause()
	Stop()
	Load(data []byte) error
	IsPlaying() bool
	GetSongInfo() string
}

type Iterator interface {
	Next() ([]byte, error)
	Pre() ([]byte, error)
}

type mp3Player struct {
	iter				Iterator
	song    			*mp3.Decoder
	stop    			chan byte
	ctx     			*oto.Context
	playing, paused 	bool
	waiting 			sync.Mutex
	startTime			int
}

func NewMp3Player(ch chan byte, iter Iterator) (Player, error) {
	otoCtx, readyChan, err := oto.NewContext(44100, 2, 2)
	//binary.Read(r, binary.BigEndian)
	if err != nil {
		return nil, err
	}
	<-readyChan
	return &mp3Player{ctx: otoCtx, iter: iter, paused: true, stop: ch, playing: false, waiting: sync.Mutex{}}, nil
}

func (m *mp3Player) Play() {
	if m.playing {
		m.stop <- 3
	} else {
		go m.play()
	}
}

func (m *mp3Player) play() {
	m.waiting.Lock()
	player := m.ctx.NewPlayer(m.song)

	defer player.Close()
	player.(io.Seeker).Seek(0, io.SeekStart)
	player.Play()
	m.startTime = int(time.Now().Unix())
	m.playing = true
	m.paused = true
	for player.IsPlaying() {
		select {
		case sig := <- m.stop:
			if sig == 1 {
				player.Pause()
				m.playing = false
				m.waiting.Unlock()
				return
			} else if sig == 2 {
				m.paused = true
				player.Pause()
				for {
					sig = <-m.stop
					if sig == 3 {
						m.paused = false
						player.Play()
						break
					} else {
						m.playing = false
						m.waiting.Unlock()
						return
					}
				}
			}
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
	player.Pause()
	m.playing = false
	m.waiting.Unlock()
	b, err := m.iter.Next()
	if err != nil {
		m.paused = true
		return
	}
	m.Load(b)
	return
}

func (m *mp3Player) Pause() {
	after := time.After(time.Second * 1)
	select {
	case m.stop <- 2:
	case <- after:
		return
	}
}

func (m *mp3Player) Stop() {
	after := time.After(time.Second * 1)
	select {
	case m.stop <- 1:
	case <- after:
		return
	}
}

func (m *mp3Player) Load(data []byte) error {
	paused := m.paused
	if m.playing {
		m.Stop()
	}

	decorded, err := mp3.NewDecoder(bytes.NewReader(data))
	if err != nil {
		return err
	}
	m.song = decorded
	if !paused {
		go m.play()
	}
	return nil
}

func (m *mp3Player) GetSongInfo() string {
	if m.playing == false {
		return "No music playing right now"
	}
	now := int(time.Now().Unix()) - m.startTime
	t := int(m.song.Length()) / 4 / 44100
	return fmt.Sprintf("#>%02d:%02d | %02d:%02d<#", t / 60, t % 60, now / 60, now % 60)
}

func (m *mp3Player) IsPlaying() bool {
	return m.playing
}

/*func (m *mp3Player) NextSong() {
	if m.songs.next != nil {
		select {
		case m.stop <- 1:
			m.songs = m.songs.next
			go m.play()
		default:
			return
		}
	} else {
		select {
		case m.stop <- 1:
			m.songs = m.first
			go m.play()
		default:
			return
		}
	}
}

func (m mp3Player) PreSong() {
	if m.songs.pre != nil {
		select {
		case m.stop <- 1:
			m.songs = m.songs.pre
			go m.play()
		default:
			return
		}
	} else {
		select {
		case m.stop <- 1:
			m.songs = m.last
			go m.play()
		default:
			return
		}
	}
}*/
