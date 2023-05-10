package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultWidth = 20
 	listHeight = 14
)

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

type model struct {
	songList  			[]string
	list       			*list.Model
	inputs     			textinput.Model
	controller 			Controller
	helpForInput		string
	info	   			string
	showInfoList		bool

	position   			int
	choice     			string
	quitting   			bool
	title	   			string

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
	//help.New()
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

func RunModel(controller Controller, s []string) {
	items := make([]list.Item, len(s))
	for i := 0; i < len(s); i++ {
		items[i] = item(s[i])
	}
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)

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

	m := model{list: &l, inputs: input, controller: controller}
	m.list.Title = "...Info..."

	tea.NewProgram(&m).Run()
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
			if len(m.songsAll[m.localRemote]) >= m.songListSixXPosition * 6 {
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
					if len(m.songListShoving) != 0 {
						err := m.controller.AddSong(m.songListShoving[m.songListPosition - 1])
						m.info = m.songListShoving[m.songListPosition - 1]
						if err == nil {
							m.info = fmt.Sprintf("Added %s", m.songListShoving[m.songListPosition - 1])
						} else {
							m.info = fmt.Sprintf("Error: %s", err.Error())
						}
					} else {
						m.info = "No song to add"
					}
				} else if m.action == 2 {
					if len(m.songListShoving) != 0 {
						err := m.controller.DeleteSong((6 * m.songListSixXPosition) + m.songListPosition)
						if err == nil {
							m.info = fmt.Sprintf("Deleted %s", m.songListShoving[m.songListPosition - 1])
						} else {
							m.info = fmt.Sprintf("Error: %s", err.Error())
						}
					 } else {
						err := m.controller.DeleteSong(0)
						if err == nil {
							m.info = fmt.Sprintf("Deleted Deleted all songs from list")
						} else {
							m.info = fmt.Sprintf("Error: %s", err.Error())
						}
					 }
				} else if m.action == 3 {
					//Не работает?
					if len(m.songListShoving) != 0 {
						err := m.controller.DeleteLocal(m.songListShoving[m.songListPosition - 1])
						if err == nil {
							m.info = fmt.Sprintf("Deleted local %s", m.songListShoving[m.songListPosition - 1])
						} else {
							m.info = fmt.Sprintf("Error: %s", err.Error())
						}
					} else {
						m.info = "No song to delete"
					}
				} else if m.action == 4 {
					if len(m.songListShoving) != 0 {
						err := m.controller.SaveSong(m.songListShoving[m.songListPosition - 1])
						if err == nil {
							m.info = fmt.Sprintf("Saved local %s", m.songListShoving[m.songListPosition - 1])
						} else {
							m.info = fmt.Sprintf("Error: %s", err.Error())
						}
					} else {
						m.info = "No song to save"
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
					if len(m.songListSearching) >= (m.songListSixXPosition + 1) * 6 + 1 {
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
					m.list.Title = m.controller.GetCurrent()
				case 1:
					m.controller.PlaySong("")
					m.info = fmt.Sprintf("Paused %s", m.controller.GetCurrent())
					m.list.Title = m.controller.GetCurrent()
				case 2:
					m.inputs.SetValue("")
					m.info = fmt.Sprintf("Paused %s", m.controller.GetCurrent())
					m.controller.PauseSong()
					m.list.Title = m.controller.GetCurrent()
				case 3:
					m.inputs.SetValue("")
					m.helpForInput = "Volume count (between 0 and 100)"
					m.showInput = true
					m.inputs.Focus()
					m.list.Title = m.controller.GetCurrent()
				case 4:
					m.inputs.SetValue("")
					err := m.controller.NextSong()
					if err == nil {
						m.info = fmt.Sprintf("")
					} else {
						m.info = fmt.Sprintf("Error: %s", err.Error())
					}
					m.inputs.Focus()
					m.list.Title = m.controller.GetCurrent()
				case 5:
					m.inputs.SetValue("")
					err := m.controller.PreSong()
					if err == nil {
						m.info = fmt.Sprintf("")
					} else {
						m.info = fmt.Sprintf("Error: %s", err.Error())
					}
					m.inputs.Focus()
					m.list.Title = m.controller.GetCurrent()
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
						m.list.Title = m.controller.GetCurrent()
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
						m.list.Title = m.controller.GetCurrent()
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
					m.list.Title = m.controller.GetCurrent()
				case 9:
					m.inputs.SetValue("")
					m.info = fmt.Sprintf("Stopped %s", m.controller.GetCurrent())
					m.controller.StopSong()
					m.inputs.Focus()
					m.list.Title = m.controller.GetCurrent()
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
					m.list.Title = m.controller.GetCurrent()
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
					m.list.Title = m.controller.GetCurrent()
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
			for i, l := range m.songListShoving {
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
		`,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7")).Width(40).Render(m.helpForInput),
		m.inputs.View(),
		s.String(),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#767676")).Width(30).Render("Press Enter for continue ->"),)
	}
	return m.info + "\n" + m.list.View()
}