package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/list"
)

var dataFileName = "tuido.json"

func getDataFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(configDir, dataFileName), nil
}

func read() ([]list.Item, error) {
	var todos []todo

	filePath, err := getDataFilePath()
	if err != nil {
		return nil, err
	}
	fp, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	file, err := io.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	// if file is empty fallback to empty json array
	if len(file) == 0 {
		file = []byte("[]")
	}
	err = json.Unmarshal(file, &todos)
	if err != nil {
		return nil, err
	}

	var items = make([]list.Item, len(todos))
	for i, t := range todos {
		items[i] = t
	}
	return items, nil
}

func write(items []list.Item) error {
	var todos = make([]todo, len(items))
	for i, maybeTodo := range items {
		t, ok := maybeTodo.(todo)
		if !ok {
			return errors.New("list.Item does not cast to todo struct.")
		}
		todos[i] = t
	}

	file, err := json.Marshal(todos)
	if err != nil {
		return err
	}
	filePath, err := getDataFilePath()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, file, 0644)
	if err != nil {
		return err
	}
	return nil
}
