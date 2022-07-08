package main

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	insertItem key.Binding
	edit       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add todo"),
		),
		edit: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "edit item"),
		),
	}
}

type delegateKeyMap struct {
	toggleDone key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		toggleDone: key.NewBinding(
			key.WithKeys(" ", "enter"),
			key.WithHelp("space/enter", "toggle item"),
		),
	}
}
