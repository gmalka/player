package MusicPlayer

import (
	"io"
	"sync"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

type player interface {
	Play()
	Pause()
	Stop()
	Load(decorded *mp3.Decoder)
}

type mp3Player struct {
	song				*mp3.Decoder
	stop               	chan byte
	ctx                	*oto.Context
	playing            	bool
	waiting				sync.Mutex
}

func NewMp3Player(ch chan byte) (player, error) {
	otoCtx, readyChan, err := oto.NewContext(44100, 2, 2)
	if err != nil {
		return nil, err
	}
	<-readyChan
	return &mp3Player{ctx: otoCtx, stop: ch, playing: false, waiting: sync.Mutex{}}, nil
}

func (m *mp3Player) Play() {
	if m.playing {
		m.stop <- 3
	} else {
		m.play()
	}
}

func (m *mp3Player) play() {
	m.waiting.Lock()
	player := m.ctx.NewPlayer(m.song)

	//go func (player oto.Player)  {
		defer player.Close()
		player.(io.Seeker).Seek(0, io.SeekStart)
		player.Play()
		for player.IsPlaying() {
			select {
				case sig := <- m.stop:
					if sig == 1 {
						player.Pause()
						m.playing = false
						m.waiting.Unlock()
						return
					} else if sig == 2 {
						player.Pause()
						for {
							sig = <- m.stop
							if sig == 3 {
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
	//}(player)
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

func (m *mp3Player) Load(decorded *mp3.Decoder) {
	m.Stop()
	m.song = decorded
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

