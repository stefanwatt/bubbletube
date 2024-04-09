package config

import (
	"os"
	"path"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	CONFIG_DIR             string
	TOKEN_PATH             string
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

func InitConfig() error {
	var err error
	CONFIG_DIR, err = os.UserConfigDir()
	if err != nil {
		return err
	}
	CONFIG_DIR = path.Join(CONFIG_DIR, "bubbletube")
	if _, err := os.Stat(CONFIG_DIR); os.IsNotExist(err) {
		err = os.Mkdir(CONFIG_DIR, 0700)
		if err != nil {
			return err
		}
	}
	TOKEN_PATH = path.Join(CONFIG_DIR, "oauth-token.json")
	if _, err := os.Stat(TOKEN_PATH); os.IsNotExist(err) {
		_, err = os.Create(TOKEN_PATH)
		if err != nil {
			return err
		}
	}
	return nil
}
