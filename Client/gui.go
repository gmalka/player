package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

type model struct {
	list      *list.Model
	inputs    textinput.Model
	showInput bool
	position  int
	choice    string
	quitting  bool
	title     string
}

var i int

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "down", "j":
			if m.showInput == false && m.position < 11 {
				m.position = m.list.Index()
			}
		case "up", "k":
			if m.showInput == false && m.position > 0 {
				m.position = m.list.Index()
			}
		case "left", "h":
			if m.showInput == false {
				m.position = m.list.Index()
				/*if m.position - 7 >= 0 {
					m.position -= 7
				} else {
					m.position = 0
				}*/
			}
		case "right", "l":
			if m.showInput == false {
				m.position = m.list.Index()
				// if m.position + 7 <= 11 {
				// 	m.position += 7
				// } else {
				// 	m.position = 11
				// }
			}
		case "enter":
			switch m.position {
			case 0:
				if m.showInput == true {
					m.showInput = false
				} else {
					m.inputs.SetValue("")
					m.showInput = true
					m.inputs.Focus()
				}
			case 1:
				if m.showInput == true {
					m.showInput = false
				} else {
					m.inputs.SetValue("")
					m.showInput = true
					m.inputs.Focus()
				}
			case 2:
				if m.showInput == true {
					m.showInput = false
				} else {
					m.inputs.SetValue("")
					m.showInput = true
					m.inputs.Focus()
				}
			case 3:
				if m.showInput == true {
					m.showInput = false
				} else {
					m.inputs.SetValue("")
					m.showInput = true
					m.inputs.Focus()
				}
			case 4:
				if m.showInput == true {
					m.showInput = false
				} else {
					m.inputs.SetValue("")
					m.showInput = true
					m.inputs.Focus()
				}
			case 5:
				if m.showInput == true {
					m.showInput = false
				} else {
					m.inputs.SetValue("")
					m.showInput = true
					m.inputs.Focus()
				}
			case 6:
				if m.showInput == true {
					m.showInput = false
				} else {
					m.inputs.SetValue("")
					m.showInput = true
					m.inputs.Focus()
				}
			case 7:
				if m.showInput == true {
					m.showInput = false
				} else {
					m.inputs.SetValue("")
					m.showInput = true
					m.inputs.Focus()
				}
			case 8:
				if m.showInput == true {
					m.showInput = false
				} else {
					m.inputs.SetValue("")
					m.showInput = true
					m.inputs.Focus()
				}
			case 9:
				if m.showInput == true {
					m.showInput = false
				} else {
					m.inputs.SetValue("")
					m.showInput = true
					m.inputs.Focus()
				}
			case 10:
				if m.showInput == true {
					m.showInput = false
				} else {
					m.inputs.SetValue("")
					m.showInput = true
					m.inputs.Focus()
				}
			case 11:
				if m.showInput == true {
					m.showInput = false
				} else {
					m.inputs.SetValue("")
					m.showInput = true
					m.inputs.Focus()
				}
			}
			return m, nil
		}
	}
	if m.showInput == true {
		var cmds tea.Cmd
		m.inputs, cmds = m.inputs.Update(msg)
		return m, cmds
	}

	var cmd tea.Cmd
	m.list.Title = m.list.Title + strconv.Itoa(m.position)
	*m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.showInput == true {
		return fmt.Sprintf(`
   %s
   %s

   %s
		`,
			lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7")).Width(30).Render("Name:"),
			m.inputs.View(),
			lipgloss.NewStyle().Foreground(lipgloss.Color("#767676")).Width(30).Render("Press Enter for continue ->"))
	}
	return "\n" + m.list.View()
}

func main() {
	items := []list.Item{
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

	m := model{list: &l, inputs: input}
	m.list.Title = "...Info..."

	tea.NewProgram(&m).Run()
}
