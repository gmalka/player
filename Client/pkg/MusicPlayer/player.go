package MusicPlayer

import (
	"bytes"
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
}

type mp3Player struct {
	song    *mp3.Decoder
	stop    chan byte
	ctx     *oto.Context
	playing bool
	paused	bool
	waiting sync.Mutex
}

func NewMp3Player(ch chan byte) (Player, error) {
	otoCtx, readyChan, err := oto.NewContext(44100, 2, 2)
	if err != nil {
		return nil, err
	}
	<-readyChan
	return &mp3Player{ctx: otoCtx, paused: true, stop: ch, playing: false, waiting: sync.Mutex{}}, nil
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
	m.playing = true
	m.paused = false
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
	return
}

func (m *mp3Player) Pause() {
	select {
	case m.stop <- 2:
	default:
		return
	}
}

func (m *mp3Player) Stop() {
	select {
	case m.stop <- 1:
	default:
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
