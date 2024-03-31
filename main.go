package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	defaultWidth      = 100
	defaultHeight     = 14
	program           *tea.Program
)

func main() {
	err := godotenv.Load() // Load .env file from the current directory
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// if len(os.Args) > 1 && os.Args[1] == "auth" {
	// 	if err := Authenticate(); err != nil {
	// 		log.Fatalf("Failed to authenticate: %v", err)
	// 	}
	// 	return
	// }

	go InitMpvConn()
	if err := Authenticate(); err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	l := MapPlaylistModel()
	m := model{list: l}
	program = tea.NewProgram(m)
	if _, err = program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
