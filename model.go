package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type todo struct {
	text string
	done bool
}

func (i todo) FilterValue() string { return "" }

type editing struct {
	// editing field must be set to true to render proper editing view
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

func newModel() model {
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

	return model{list: l, textInput: ti, keys: listKeys}
}
