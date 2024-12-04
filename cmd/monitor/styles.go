package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	tableStyle = table.Styles{
		Header: lipgloss.NewStyle().Padding(0, 1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")),
		Cell: lipgloss.NewStyle().Padding(0, 1),
	}

	highlight = lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57"))

	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
)
