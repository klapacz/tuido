package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type state int

const (
	all state = iota
	new
)

type todo struct {
	text      string
	completed bool
}

func (t *todo) renderCheckbox() string {
	if t.completed {
		return "[X]"
	}
	return "[ ]"
}

func (i todo) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	t, ok := listItem.(todo)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s %s", t.renderCheckbox(), t.text)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type model struct {
	list      list.Model
	textInput textinput.Model
	state     state
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		keypress := msg.String()

		switch m.state {
		case all:
			switch keypress {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "n":
				m.state = new
				m.textInput.Reset()
				return m, cmd

			case " ":
				t, ok := m.list.SelectedItem().(todo)
				if !ok {
					return m, nil
				}
				t.completed = !t.completed
				return m, m.list.SetItem(m.list.Index(), t)
			}
		case new:
			switch keypress {
			case "enter":
				if val := m.textInput.Value(); val != "" {
					cmd = m.list.InsertItem(0, todo{val, false})
				}
				m.state = all
				return m, cmd
			case "esc":
				m.state = all
				return m, nil
			}
		}
	}

	switch m.state {
	case all:
		m.list, cmd = m.list.Update(msg)
	case new:
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.state == all {
		return "\n" + m.list.View()
	}
	return "\n" + m.textInput.View()
}

func main() {
	items := []list.Item{
		todo{"first", false},
		todo{"second", false},
		todo{"third", true},
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "todos"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	ti := textinput.New()
	ti.Placeholder = "New todo text"
	ti.CharLimit = 156
	ti.Width = 20
	ti.Focus()

	m := model{list: l, textInput: ti}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
