package Player

import (
	"sync"
	"io"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)


type song struct {
	s         	*mp3.Decoder
	next, pre 	*song
}

type Player interface {
	Play()
	Pause()
	NextSong()
	PreSong()
	Add(decorded *mp3.Decoder)
}

type mp3Player struct {
	songs, last, first	*song
	stop               	chan byte
	ctx                	*oto.Context
	playing            	bool
	waiting			   	*sync.Mutex
}

func NewMp3Player(ch chan byte) (Player, error) {
	otoCtx, readyChan, err := oto.NewContext(44100, 2, 2)
	if err != nil {
		return nil, err
	}
	<-readyChan
	return &mp3Player{ctx: otoCtx, stop: ch, playing: false, waiting: &sync.Mutex{}}, nil
}

func (m *mp3Player) Play() {
	if m.playing {
		m.stop <- 3
	} else {
		m.play()
	}
}

func (m *mp3Player) play() {
	if m.songs == nil {
		return
	}
	m.waiting.Lock()
	player := m.ctx.NewPlayer(m.songs.s)
	defer player.Close()
	player.(io.Seeker).Seek(0, io.SeekStart)
	player.Play()
	m.playing = true

	go func (player oto.Player)  {
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
		m.NextSong()
		m.playing = false
		m.waiting.Unlock()
		m.Play()
	}(player)
}

func (m *mp3Player) Pause() {
	select {
		case m.stop <- 2:
	default:
		return
	}
}

func (m *mp3Player) NextSong() {
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
}

func (m *mp3Player) Add(decorded *mp3.Decoder) {
	if m.songs == nil {
		m.songs = &song{s: decorded}
		m.last = m.songs
		m.first = m.songs
	
	} else {
		m.last.next = &song{s: decorded, pre: m.last}
		m.last = m.last.next
	}
}