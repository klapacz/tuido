package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type itemDelegate struct {
	list.DefaultDelegate
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	t, ok := listItem.(todo)
	if !ok {
		return
	}

	c := "[ ]"
	if t.Done {
		c = "[X]"
	}
	str := fmt.Sprintf("%s %s", c, t.Text)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

func newItemDelegate(keys *delegateKeyMap) itemDelegate {
	d := itemDelegate{}

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		t, ok := m.SelectedItem().(todo)
		if !ok {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.toggle):
				t.Done = !t.Done
				return m.SetItem(m.Index(), t)
			}
		}
		return nil
	}

	help := []key.Binding{keys.toggle}
	d.ShortHelpFunc = func() []key.Binding {
		return help
	}
	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}
