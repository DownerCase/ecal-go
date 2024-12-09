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
