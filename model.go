package main

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type todo struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
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
	todos, err := read()
	if err != nil {
		log.Fatalf("Error while reading data file: %s", err)
	}

	var (
		listKeys     = newListKeyMap()
		delegateKeys = newDelegateKeyMap()
	)

	l := list.New(todos, newItemDelegate(delegateKeys), 0, 0)
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
