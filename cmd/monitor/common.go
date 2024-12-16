package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func NewTable(columns []table.Column) table.Model {
	return table.New(
		table.WithHeight(8),
		table.WithFocused(true),
		table.WithStyles(tableStyle),
		table.WithColumns(columns),
	)
}

type NavKeyMap map[string]func() tea.Cmd

func (navKeys NavKeyMap) HandleMsg(msg tea.Msg) (cmd tea.Cmd, navigated bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if f, ok := navKeys[msg.String()]; ok {
			return f(), true
		}
	}
	return nil, false
}
