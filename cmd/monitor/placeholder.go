package main

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

func placeholderTab(page Page) string {
	style := lipgloss.NewStyle().Height(12)
	return style.Render("Page: ", strconv.FormatInt(int64(page+1), 10))
}
