package config

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	DefaultWidth       = 100
	DefaultHeight      = 10
	DefaultVolume      = 50.0
	DefaultVolumeWidth = 20
	TitleStyle         = lipgloss.NewStyle().MarginLeft(2)
	ItemStyle          = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle  = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	PaginationStyle    = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle          = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	QuitTextStyle      = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)
