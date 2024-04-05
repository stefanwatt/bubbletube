package main

import (
	"fmt"
	"log"
	"os"

	controller "bubbletube/controller"
	model "bubbletube/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

var program *tea.Program

func main() {
	err := godotenv.Load() // Load .env file from the current directory
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := model.Authenticate(); err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	playlists := model.MapPlaylistsModel(controller.PlaylistDelegate{})
	screen := model.MapDefaultScreen(playlists, controller.SongDelegate{})
	sc := controller.NewScreenController(screen)

	program = tea.NewProgram(sc)
	go model.InitMpvConn(program)
	if _, err = program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
