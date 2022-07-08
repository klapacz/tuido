package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle          = lipgloss.NewStyle().Padding(1, 2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type todo struct {
	text string
	done bool
}

func (i todo) FilterValue() string { return "" }

type editing struct {
	editing bool
	index   int
	done    bool
}

type view struct {
	adding  bool
	editing editing
}

type model struct {
	list      list.Model
	textInput textinput.Model
	keys      *listKeyMap
	view      view
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func updateList(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.insertItem):
			m.textInput.Placeholder = "New todo text"
			m.view = view{adding: true}
			return m, nil
		case key.Matches(msg, m.keys.edit):
			if t, ok := m.list.SelectedItem().(todo); ok {
				m.textInput.Placeholder = ""
				m.textInput.SetValue(t.text)
				e := editing{done: t.done, index: m.list.Index(), editing: true}
				m.view = view{editing: e}
			}
			return m, nil
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func updateAdding(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keypress := msg.String()
		switch keypress {
		case "enter":
			if val := m.textInput.Value(); val != "" {
				m.textInput.Reset()
				if m.view.adding {
					cmd = m.list.InsertItem(0, todo{val, false})
				} else {
					t := todo{val, m.view.editing.done}
					cmd = m.list.SetItem(m.view.editing.index, t)
				}
			}
			m.view = view{}
			return m, cmd
		case "esc":
			m.view = view{}
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.view.adding || m.view.editing.editing {
		return updateAdding(msg, m)
	}
	return updateList(msg, m)
}

func (m model) View() string {
	if m.view.adding || m.view.editing.editing {
		return appStyle.Render(m.textInput.View())
	}
	return appStyle.Render(m.list.View())
}

func main() {
	items := []list.Item{
		todo{"first", false},
		todo{"second", false},
		todo{"third", true},
	}

	var (
		listKeys     = newListKeyMap()
		delegateKeys = newDelegateKeyMap()
	)

	l := list.New(items, newItemDelegate(delegateKeys), 0, 0)
	l.Title = "tuido"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.insertItem,
			listKeys.edit,
		}
	}

	ti := textinput.New()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Focus()

	m := model{list: l, textInput: ti, keys: listKeys}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
