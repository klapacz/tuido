package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle          = lipgloss.NewStyle().Padding(1, 2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			err := write(m.list.Items())
			if err != nil {
				log.Fatalf("Error while writing data file: %s", err)
			}
			return m, tea.Quit
		}
	}

	if m.view.adding || m.view.editing.editing {
		return updateTextInput(msg, m)
	}
	return updateList(msg, m)
}

func updateList(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case startEditingMsg:
		m.input.Placeholder = ""
		m.input.SetValue(msg.todo.Text)
		e := editing{done: msg.todo.Done, index: msg.index, editing: true}
		m.view = view{editing: e}
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.insert):
			m.input.Placeholder = "New todo text"
			m.view = view{adding: true}
			return m, nil
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func updateTextInput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keypress := msg.String()
		switch keypress {
		case "enter":
			if val := m.input.Value(); val != "" {
				m.input.Reset()
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

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.view.adding || m.view.editing.editing {
		return appStyle.Render(m.input.View())
	}
	return appStyle.Render(m.list.View())
}

func main() {
	m := newModel()

	if err := tea.NewProgram(m).Start(); err != nil {
		log.Println("Error running program:", err)
		os.Exit(1)
	}
}
