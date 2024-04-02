package config

import "github.com/charmbracelet/lipgloss"

type ColorPalette struct {
	Rosewater lipgloss.Color
	Flamingo  lipgloss.Color
	Pink      lipgloss.Color
	Mauve     lipgloss.Color
	Red       lipgloss.Color
	Maroon    lipgloss.Color
	Peach     lipgloss.Color
	Yellow    lipgloss.Color
	Green     lipgloss.Color
	Teal      lipgloss.Color
	Sky       lipgloss.Color
	Sapphire  lipgloss.Color
	Blue      lipgloss.Color
	Lavender  lipgloss.Color
	Text      lipgloss.Color
	Subtext1  lipgloss.Color
	Subtext0  lipgloss.Color
	Overlay2  lipgloss.Color
	Overlay1  lipgloss.Color
	Overlay0  lipgloss.Color
	Surface2  lipgloss.Color
	Surface1  lipgloss.Color
	Surface0  lipgloss.Color
	Base      lipgloss.Color
	Mantle    lipgloss.Color
	Crust     lipgloss.Color
}

var Colors = ColorPalette{
	Rosewater: lipgloss.Color("#f2d5cf"),
	Flamingo:  lipgloss.Color("#eebebe"),
	Pink:      lipgloss.Color("#f4b8e4"),
	Mauve:     lipgloss.Color("#ca9ee6"),
	Red:       lipgloss.Color("#e78284"),
	Maroon:    lipgloss.Color("#ea999c"),
	Peach:     lipgloss.Color("#ef9f76"),
	Yellow:    lipgloss.Color("#e5c890"),
	Green:     lipgloss.Color("#a6d189"),
	Teal:      lipgloss.Color("#81c8be"),
	Sky:       lipgloss.Color("#99d1db"),
	Sapphire:  lipgloss.Color("#85c1dc"),
	Blue:      lipgloss.Color("#8caaee"),
	Lavender:  lipgloss.Color("#babbf1"),
	Text:      lipgloss.Color("#c6d0f5"),
	Subtext1:  lipgloss.Color("#b5bfe2"),
	Subtext0:  lipgloss.Color("#a5adce"),
	Overlay2:  lipgloss.Color("#949cbb"),
	Overlay1:  lipgloss.Color("#838ba7"),
	Overlay0:  lipgloss.Color("#737994"),
	Surface2:  lipgloss.Color("#626880"),
	Surface1:  lipgloss.Color("#51576d"),
	Surface0:  lipgloss.Color("#414559"),
	Base:      lipgloss.Color("#303446"),
	Mantle:    lipgloss.Color("#292c3c"),
	Crust:     lipgloss.Color("#232634"),
}
