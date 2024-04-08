package config

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	DefaultWidth           = 100
	DefaultQueuePanelWidth = 45
	DefaultHeight          = 10
	DefaultWaitlistHeight  = 12
	DefaultVolume          = 50.0
	DefaultVolumeWidth     = 20
	TitleStyle             = lipgloss.NewStyle().
				PaddingLeft(2).
				PaddingRight(2).
				Align(lipgloss.Left).
				Background(Colors.Text).
				Foreground(Colors.Crust)
	PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	QuitTextStyle   = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)
