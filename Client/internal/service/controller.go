package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type model struct {
	songList  		[]string
	list       		*list.Model
	inputs     		textinput.Model
	controller 		*myController
	helpForInput	string
	info	   		string
	showInfoList	bool

	position   		int
	choice     		string
	quitting   		bool
	title	   		string

	songsAll			[][]string
	localRemote			int
	showAll				bool

	showList  			 bool
	showInput			 bool
	songListShoving 	 []string
	songListSearching 	 []string
	songListAll			 []string
	preCommand			 string
	songListPosition	 int
	songListSixXPosition int

	action				 int
}

var i int

type myController struct {
	player player
	songmanager songmanager
}

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	help.New()
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " ") + " <")
		}
	}

	_ = fn
	_ = str
	fmt.Fprint(w, fn(str))
}

func (m *model) Init() tea.Cmd {
	return  nil
}

func showAll(m *model, msg tea.Msg) *model {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "enter", tea.KeyEsc.String():
			m.showAll = false
			m.showInput = false
			return m
		case "down":
			m.songListSixXPosition++
			if len(m.songsAll[m.localRemote]) > m.songListSixXPosition * 6 {
				i := 0
				for ; i < len(m.songListShoving) && m.songListSixXPosition * 6 + i < len(m.songsAll[m.localRemote]); i++ {
					m.songListShoving[i] = m.songsAll[m.localRemote][m.songListSixXPosition * 6 + i]
				}
				for ; i < 6 && m.songListSixXPosition * 6 + i < len(m.songsAll[m.localRemote]); i++ {
					m.songListShoving = append(m.songListShoving, m.songsAll[m.localRemote][m.songListSixXPosition * 6 + i])
				}
				m.songListShoving = m.songListShoving[:i]
			} else if m.localRemote == 0 {
				m.localRemote++
				m.songListSixXPosition = 0
				i := 0
				for ; i < len(m.songListShoving) && i < len(m.songsAll[m.localRemote]); i++ {
					m.songListShoving[i] = m.songsAll[m.localRemote][i]
				}
				for ; i < 6 && i < len(m.songsAll[m.localRemote]); i++ {
					m.songListShoving = append(m.songListShoving, m.songsAll[m.localRemote][i])
				}
				m.songListShoving = m.songListShoving[:i]
			} else {
				m.songListSixXPosition--
			}
		case "up":
			m.songListSixXPosition--
			if m.songListSixXPosition >= 0 {
				i := 0
				for ; i < len(m.songListShoving) && m.songListSixXPosition * 6 + i < len(m.songsAll[m.localRemote]); i++ {
					m.songListShoving[i] = m.songsAll[m.localRemote][m.songListSixXPosition * 6 + i]
				}
				for ; i < 6 && m.songListSixXPosition * 6 + i < len(m.songsAll[m.localRemote]); i++ {
					m.songListShoving = append(m.songListShoving, m.songsAll[m.localRemote][m.songListSixXPosition * 6 + i])
				}
				m.songListShoving = m.songListShoving[:i]
			} else if m.localRemote == 1 {
				m.localRemote--
				m.songListSixXPosition = len(m.songsAll[m.localRemote]) / 6
				i := 0
				for ; i < len(m.songListShoving) && m.songListSixXPosition * 6 + i < len(m.songsAll[m.localRemote]); i++ {
					m.songListShoving[i] = m.songsAll[m.localRemote][m.songListSixXPosition * 6 + i]
				}
				for ; i < 6 && m.songListSixXPosition * 6 + i < len(m.songsAll[m.localRemote]); i++ {
					m.songListShoving = append(m.songListShoving, m.songsAll[m.localRemote][m.songListSixXPosition * 6 + i])
				}
				m.songListShoving = m.songListShoving[:i]
			} else {
				m.songListSixXPosition++
			}
		}
	}
	return nil
}

func showInfoList(m *model, msg tea.Msg) *model {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "enter", tea.KeyEsc.String():
			m.showInfoList = false
			m.showInput = false
			return m
		case "down":
			m.songListSixXPosition++
			if m.songListSixXPosition * 10 < len(m.songList) {
				i := 0
				for ; i < len(m.songListShoving) && m.songListSixXPosition * 10 + i < len(m.songList); i++ {
					m.songListShoving[i] = m.songList[m.songListSixXPosition * 10 + i]
				}
				for ; i < 10 && m.songListSixXPosition * 10 + i < len(m.songList); i++ {
					m.songListShoving = append(m.songListShoving, m.songList[m.songListSixXPosition * 10 + i])
				}
				m.songListShoving = m.songListShoving[:i]
			} else {
				m.songListSixXPosition--
			}
		case "up":
			m.songListSixXPosition--
			if m.songListSixXPosition >= 0 {
				i := 0
				for ; i < len(m.songListShoving) && m.songListSixXPosition * 10 + i < len(m.songList); i++ {
					m.songListShoving[i] = m.songList[m.songListSixXPosition * 10 + i]
				}
				for ; i < 10 && m.songListSixXPosition * 10 + i < len(m.songList); i++ {
					m.songListShoving = append(m.songListShoving, m.songList[m.songListSixXPosition * 10 + i])
				}
				m.songListShoving = m.songListShoving[:i]
			} else {
				m.songListSixXPosition = 0
			}
		}
	}
	return nil
}

func showInput(m *model, msg tea.Msg) *model {
	if m.showList {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.list.SetWidth(msg.Width)
			return m
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "ctrl+c", tea.KeyEsc.String(), "q":
				m.showList = false
				m.showInput = false
				m.songListPosition = 0
				m.songListSixXPosition = 0
				return m
			case "enter":
				if m.action == 1 {
					err := m.controller.AddSong(m.songListShoving[m.songListPosition - 1])
					m.info = m.songListShoving[m.songListPosition - 1]
					if err == nil {
						m.info = fmt.Sprintf("Added %s", m.songListShoving[m.songListPosition - 1])
					} else {
						m.info = fmt.Sprintf("Error: %s", err.Error())
					}
				} else if m.action == 2 {
					err := m.controller.DeleteSong(m.songListShoving[m.songListPosition - 1])
					if err == nil {
						m.info = fmt.Sprintf("Deleted %s", m.songListShoving[m.songListPosition - 1])
					} else {
						m.info = fmt.Sprintf("Error: %s", err.Error())
					}
				} else if m.action == 3 {

					//Не работает?
					err := m.controller.DeleteLocal(m.songListShoving[m.songListPosition - 1])
					if err == nil {
						m.info = fmt.Sprintf("Deleted local %s", m.songListShoving[m.songListPosition - 1])
					} else {
						m.info = fmt.Sprintf("Error: %s", err.Error())
					}
				} else if m.action == 4 {
					err := m.controller.SaveSong(m.songListShoving[m.songListPosition - 1])
					if err == nil {
						m.info = fmt.Sprintf("Saved local %s", m.songListShoving[m.songListPosition - 1])
					} else {
						m.info = fmt.Sprintf("Error: %s", err.Error())
					}
				}
				m.action = 0
				m.showList = false
				m.showInput = false
				m.songListPosition = 1
				m.songListSixXPosition = 0
				return m
			case "down":
				m.songListPosition++
				if m.songListPosition > len(m.songListShoving) {
					if len(m.songListSearching) > (m.songListSixXPosition + 1) * 6 + 1 {
						m.songListSixXPosition++
						i := 0
						for ; i < len(m.songListShoving) && m.songListSixXPosition * 6 + i < len(m.songListSearching); i++ {
							m.songListShoving[i] = m.songListSearching[m.songListSixXPosition * 6 + i]
						}
						for ; i < 6 && m.songListSixXPosition * 6 + i < len(m.songListSearching); i++ {
							m.songListShoving = append(m.songListShoving, m.songListSearching[m.songListSixXPosition * 6 + i])
						}
						m.songListShoving = m.songListShoving[:i]
					}
					m.songListPosition = 1
				}
			case "up":
				m.songListPosition--
				if m.songListPosition < 1 {
					if m.songListSixXPosition > 0 {
						m.songListSixXPosition--
						i := 0
						for ; i < len(m.songListShoving) && m.songListSixXPosition * 6 + i < len(m.songListSearching); i++ {
							m.songListShoving[i] = m.songListSearching[m.songListSixXPosition * 6 + i]
						}
						for ; i < 6 && m.songListSixXPosition * 6 + i < len(m.songListSearching); i++ {
							m.songListShoving = append(m.songListShoving, m.songListSearching[m.songListSixXPosition * 6 + i])
						}
						m.songListShoving = m.songListShoving[:i]
						m.songListPosition = 6
					} else {
						m.songListPosition = 1
					}
				}
			}
		}
		command := m.inputs.Value()
		if command != m.preCommand {
			s := make([]string, 0, 10)
			for _, l := range m.songListAll {
				if strings.Contains(l, command) {
					s = append(s, l)
				}
			}
			m.songListSearching = s
			i := 0
			for ; i < len(m.songListShoving) && m.songListSixXPosition * 6 + i < len(m.songListSearching); i++ {
				m.songListShoving[i] = m.songListSearching[m.songListSixXPosition * 6 + i]
			}
			for ; i < 6 && m.songListSixXPosition * 6 + i < len(m.songListSearching); i++ {
				m.songListShoving = append(m.songListShoving, m.songListSearching[m.songListSixXPosition * 6 + i])
			}
			m.songListShoving = m.songListShoving[:i]
			m.preCommand = command
		}
	} else {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.list.SetWidth(msg.Width)
			return m
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "ctrl+c", tea.KeyEsc.String(), "enter":
				m.showInput = false
				m.songListPosition = 0
				m.songListSixXPosition = 0

				command := m.inputs.Value()
				if m.list.Index() == 3 {
					m.controller.SetVolume(command)
				}
				return m
			}
		}
	}
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.showAll {
		h := showAll(m, msg)
		if h != nil {
			return m, nil
		}
	} else if m.showInfoList {
		h := showInfoList(m, msg)
		if h != nil {
			return m, nil
		}	
	} else if m.showInput {
		h := showInput(m, msg)
		if h != nil {
			return m, nil
		}	
	} else {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.list.SetWidth(msg.Width)
			return m, nil
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit
			case "down", "j":
				if m.position < 11 {
					m.position = m.list.Index()
				}
			case "up", "k":
				if m.position > 0 {
					m.position = m.list.Index()
				}
			case "left", "h":
				m.position = m.list.Index()
			case "right", "l":
				m.position = m.list.Index()
			case "enter":
				switch m.list.Index() {
				case 0:
					m.inputs.SetValue("")
					s, err := m.controller.GetAllSongs("")
					if err != nil {
						m.info = fmt.Sprintf("Error: %s", err.Error())
						break
					}
					m.songListAll = s[0]
					m.songListAll = append(m.songListAll, s[1]...)
					sort.Strings(m.songListAll)
					m.songListShoving = make([]string, 0, 6)
					i := 0
					for ; i < 6 && i < len(m.songListAll); i++ {
						m.songListShoving  = append(m.songListShoving, m.songListAll[i])
					}
					m.songListShoving = m.songListShoving[:i]
					m.songListSearching = m.songListAll
					m.helpForInput = "Song name:"
					m.showInput = true
					m.showList = true
					m.songListPosition = 1
					m.songListSixXPosition = 0
					m.action = 1
					m.list.Title = m.controller.songmanager.GetCurrent()
				case 1:
					m.controller.PlaySong("")
					m.info = fmt.Sprintf("Paused %s", m.controller.songmanager.GetCurrent())
					m.list.Title = m.controller.songmanager.GetCurrent()
				case 2:
					m.inputs.SetValue("")
					m.info = fmt.Sprintf("Paused %s", m.controller.songmanager.GetCurrent())
					m.controller.PauseSong()
					m.list.Title = m.controller.songmanager.GetCurrent()
				case 3:
					m.inputs.SetValue("")
					m.helpForInput = "Volume count (between 0 and 100)"
					m.showInput = true
					m.inputs.Focus()
					m.list.Title = m.controller.songmanager.GetCurrent()
				case 4:
					m.inputs.SetValue("")
					err := m.controller.NextSong()
					if err == nil {
						m.info = fmt.Sprintf("")
					} else {
						m.info = fmt.Sprintf("Error: %s", err.Error())
					}
					m.inputs.Focus()
					m.list.Title = m.controller.songmanager.GetCurrent()
				case 5:
					m.inputs.SetValue("")
					err := m.controller.PreSong()
					if err == nil {
						m.info = fmt.Sprintf("")
					} else {
						m.info = fmt.Sprintf("Error: %s", err.Error())
					}
					m.inputs.Focus()
					m.list.Title = m.controller.songmanager.GetCurrent()
				case 6:
					//play list
					m.inputs.SetValue("")
					s, err := m.controller.GetPlayList()
					if err != nil {
						m.info = err.Error()
					} else {
						m.songList = s
						m.songListShoving = make([]string, 0, 10)
						i := 0
						for ; i < 10 && i < len(s); i++ {
							m.songListShoving = append(m.songListShoving, s[i])
						}
						m.songListShoving = m.songListShoving[:i]
						m.showInput = true
						m.showInfoList = true
						m.list.Title = m.controller.songmanager.GetCurrent()
					}
				case 7:
					//GetAllSongs
					s, err := m.controller.GetAllSongs("")
					if err != nil {
						m.info = err.Error()
					} else {
						m.showAll = true
						sort.Strings(s[0])
						sort.Strings(s[1])
						m.songListShoving = make([]string, 0, 6)
						m.songsAll = s
						i := 0
						for ; i < 6 && i < len(s[0]); i++ {
							m.songListShoving = append(m.songListShoving, s[0][i])
						}
						m.songListShoving = m.songListShoving[:i]
						m.showInput = true
						m.showAll = true
						m.localRemote = 0
						m.songListSixXPosition = 0
						m.list.Title = m.controller.songmanager.GetCurrent()
					}
				case 8:
					//Delete
					m.inputs.SetValue("")
					s, err := m.controller.GetPlayList()
					if err != nil {
						m.info = fmt.Sprintf("Error: %s", err.Error())
						break
					}
					m.songListAll = s
					sort.Strings(m.songListAll)
					m.songListShoving = make([]string, 0, 6)
					i := 0
					for ; i < 6 && i < len(m.songListAll); i++ {
						m.songListShoving  = append(m.songListShoving, m.songListAll[i])
					}
					m.songListShoving = m.songListShoving[:i]
					m.songListSearching = m.songListAll
					m.helpForInput = "Song name:"
					m.showInput = true
					m.showList = true
					m.songListPosition = 1
					m.action = 2
					m.songListSixXPosition = 0
					m.list.Title = m.controller.songmanager.GetCurrent()
				case 9:
					m.inputs.SetValue("")
					m.info = fmt.Sprintf("Stopped %s", m.controller.songmanager.GetCurrent())
					m.controller.StopSong()
					m.inputs.Focus()
					m.list.Title = m.controller.songmanager.GetCurrent()
				case 10:
					//delete local
					m.inputs.SetValue("")
					s, err := m.controller.GetAllSongs("local")
					if err != nil {
						m.info = fmt.Sprintf("Error: %s", err.Error())
						break
					}
					m.songListAll = s[0]
					sort.Strings(m.songListAll)
					m.songListShoving = make([]string, 0, 6)
					i := 0
					for ; i < 6 && i < len(m.songListAll); i++ {
						m.songListShoving  = append(m.songListShoving, m.songListAll[i])
					}
					m.songListShoving = m.songListShoving[:i]
					m.songListSearching = m.songListAll
					m.helpForInput = "Song name:"
					m.showInput = true
					m.showList = true
					m.songListPosition = 1
					m.action = 3
					m.songListSixXPosition = 0
					m.list.Title = m.controller.songmanager.GetCurrent()
				case 11:
					//save local
					m.inputs.SetValue("")
					s, err := m.controller.GetAllSongs("remote")
					if err != nil {
						m.info = fmt.Sprintf("Error: %s", err.Error())
						break
					}
					m.songListAll = s[1]
					sort.Strings(m.songListAll)
					m.songListShoving = make([]string, 0, 6)
					i := 0
					for ; i < 6 && i < len(m.songListAll); i++ {
						m.songListShoving  = append(m.songListShoving, m.songListAll[i])
					}
					m.songListShoving = m.songListShoving[:i]
					m.songListSearching = m.songListAll
					m.helpForInput = "Song name:"
					m.showInput = true
					m.showList = true
					m.songListPosition = 1
					m.action = 4
					m.songListSixXPosition = 0
					m.list.Title = m.controller.songmanager.GetCurrent()
				}
				return m, nil
			}
		}
	}
	if m.showInput == true {
		var cmds tea.Cmd
		m.inputs, cmds = m.inputs.Update(msg)
		return m, cmds
	}

	var cmd tea.Cmd
	*m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	// if m.showLocalRemote {
	// 	s := strings.Builder{}
	// 	if m.localRemote[0] != nil {
	// 		s.WriteString("\n->Local:\n")
	// 		for _, l := range m.localRemote[0] {
	// 			s.WriteString(l)
	// 			s.WriteRune('\n')
	// 		}
	// 	}
	// 	if m.localRemote[1] != nil {
	// 		s.WriteString("\n->Remote:\n")
	// 		for _, l := range m.localRemote[1] {
	// 			s.WriteString(l)
	// 			s.WriteRune('\n')
	// 		}
	// 	}
	// 	return s.String()
	// }
	if m.showAll == true {
		s := strings.Builder{}
		if m.localRemote == 0 {
			s.WriteString("Local: \n")
		} else {
			s.WriteString("Remote: \n")
		}
		for i := 0; i < len(m.songListShoving); i++ {
			s.WriteString(fmt.Sprintf("%d)", i + 1))
			s.WriteString(m.songListShoving[i])
			s.WriteRune('\n')
		}
		if len(m.songListShoving)== 0 {
			s.WriteString("List empty\n")
		}
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#767676")).Width(30).Render("\nESC/q/Enter to back"))
		return s.String()
	}
	if m.showInput == true {
		if m.showInfoList {
			s := strings.Builder{}
			s.WriteRune('\n')
			for _, l := range m.songListShoving {
				s.WriteString(fmt.Sprintf("%d)", i + 1))
				s.WriteString(l)
				s.WriteRune('\n')
			}
			s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#767676")).Width(30).Render("\nESC/q/Enter to back"))
			return s.String()
		}
		if !m.showList {
			return fmt.Sprintf(`
   %s
   %s
   %s
		`,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7")).Width(40).Render(m.helpForInput),
		m.inputs.View(),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#767676")).Width(30).Render("Press Enter for continue ->"),)
		}
		s := strings.Builder{}
		s.WriteRune('\n')
		for i, l := range m.songListShoving {
			if i + 1 == m.songListPosition {
				s.WriteRune('>')
			}
			s.WriteString(fmt.Sprintf("%d)", i + 1))
			s.WriteString(l)
			s.WriteRune('\n')
		}
		return fmt.Sprintf(`
   %s
   %s
%s
   %s
   %d
   %d
		`,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7")).Width(40).Render(m.helpForInput),
		m.inputs.View(),
		s.String(),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#767676")).Width(30).Render("Press Enter for continue ->"),
	m.songListPosition,
	m.songListSixXPosition)
	}
	return m.info + "\n" + m.list.View()
}

func NewController(MusicPlayer player, musicFileManager songmanager) myController {
	return myController{player: MusicPlayer, songmanager: musicFileManager}
}

func (c *myController) Run() {
	items := []list.Item {
		item("Add"),
		item("Play"),
		item("Pause"),
		item("Set Volume"),
		item("Next"),
		item("Pre"),
		item("Playlist"),
		item("Get all songs"),
		item("Delete"),
		item("Stop"),
		item("Delete from local storage"),
		item("Save song in local storage"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	
	l.Title = "some some"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	var input textinput.Model

	input = textinput.New()
	input.Placeholder = "text here"
	input.Focus()
	input.Width = 30
	input.Prompt = ""

	m := model{list: &l, inputs: input, controller: c}
	m.list.Title = "...Info..."

	tea.NewProgram(&m).Run()
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

func (c *myController) DeleteSong(str string) error {
	if str != "" {
		err := c.songmanager.Delete(str)
		if err != nil {
			return err
		}
	} else {
		err := c.songmanager.Delete("")
		if err != nil {
			return err
		}
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