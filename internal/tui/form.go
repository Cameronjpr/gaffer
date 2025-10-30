package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

type menuKeyMap struct {
	Exit     key.Binding
	formKeys *huh.KeyMap
}

func defaultMenuKeyMap() *menuKeyMap {
	keys := &menuKeyMap{
		formKeys: huh.NewDefaultKeyMap(),
	}
	keys.formKeys.Quit.SetEnabled(false)
	return keys
}
