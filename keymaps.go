package main

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	insert key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insert: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add todo"),
		),
	}
}

type delegateKeyMap struct {
	toggle key.Binding
	edit   key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		toggle: key.NewBinding(
			key.WithKeys(" ", "enter"),
			key.WithHelp("space/enter", "toggle item"),
		),
		edit: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "edit item"),
		),
	}
}
