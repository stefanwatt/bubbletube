package model

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/DexterLB/mpvipc"
	tea "github.com/charmbracelet/bubbletea"
)

type MPVFloatValueChangedEvent struct {
	ID    int64
	Value float64
}

func mpvEventCmd(event mpvipc.Event) tea.Cmd {
	return func() tea.Msg {
		if event.Data == nil {
			return MPVFloatValueChangedEvent{
				ID:    event.ID,
				Value: 0,
			}
		}
		return MPVFloatValueChangedEvent{
			ID:    event.ID,
			Value: event.Data.(float64),
		}
	}
}

var (
	conn              *mpvipc.Connection
	playing           = false
	volume            = 50.0
	time_remaining    = 0.0
	time_pos          = 0.0
	currentMPVProcess *exec.Cmd
	started           = false
)

func InitMpvConn(program *tea.Program) {
	// Start MPV process with idle flag and no video
	started = true
	currentMPVProcess = exec.Command(
		"mpv",
		"--idle",
		"--input-unix-socket=/tmp/mpv_socket",
		"--no-video",
		"--volume="+fmt.Sprintf("%f", volume),
	)
	if err := currentMPVProcess.Start(); err != nil {
		fmt.Printf("Failed to start MPV: %v\n", err)
		return
	}
	time.Sleep(1 * time.Second)
	conn = mpvipc.NewConnection("/tmp/mpv_socket")
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Handle MPV events
	events, stopListening := conn.NewEventListener()
	go func() {
		for event := range events {
			msg := mpvEventCmd(*event)().(MPVFloatValueChangedEvent)
			if program != nil {
				program.Send(msg)
			}
		}
	}()

	go func() {
		conn.WaitUntilClosed()
		stopListening <- struct{}{}
	}()

	// Observe properties
	observeProperties()

	// Setup signal handling to cleanly exit MPV on interrupt
	setupSignalHandling()
}

func observeProperties() {
	properties := []string{"volume", "duration", "percent-pos", "time-pos", "time-remaining"}
	for id, prop := range properties {
		if _, err := conn.Call("observe_property", id+1, prop); err != nil {
			fmt.Println(err)
		}
	}
}

func setupSignalHandling() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for range sigChan {
			KillMpv()
		}
	}()
}

func SelectSong(item SongItem) {
	if conn == nil {
		fmt.Println("MPV not started")
		return
	}

	// Load and play the selected song
	_, err := conn.Call("loadfile", "https://www.youtube.com/watch?v="+item.VideoID, "replace")
	if err != nil {
		fmt.Printf("Failed to load song: %v\n", err)
	}
}

func TogglePlayback() bool {
	playing = !playing
	if err := conn.Set("pause", playing); err != nil {
		log.Fatal(err)
	}
	return playing
}

// VolumeUp, VolumeDown, SkipForward, SkipBackward remain unchanged
func VolumeUp() {
	if conn == nil {
		fmt.Println("MPV not started")
		return
	}
	volume = volume + 5
	err := conn.Set("volume", volume)
	if err != nil {
		fmt.Println(err)
	}
}

func VolumeDown() {
	if conn == nil {
		fmt.Println("MPV not started")
		return
	}
	volume = volume - 5
	err := conn.Set("volume", volume)
	if err != nil {
		fmt.Println(err)
	}
}

func SkipForward() {
	if conn == nil {
		fmt.Println("MPV not started")
		return
	}
	skipBy := math.Floor(math.Min(10, time_remaining))

	_, err := conn.Call("seek", skipBy, "relative", "exact")
	if err != nil {
		fmt.Println(err)
	}
}

func SkipBackward() {
	if conn == nil {
		fmt.Println("MPV not started")
		return
	}

	skipBy := -math.Floor(math.Min(10, time_pos))
	_, err := conn.Call("seek", skipBy, "relative", "exact")
	if err != nil {
		fmt.Println(err)
	}
}

func KillMpv() {
	if currentMPVProcess == nil || currentMPVProcess.Process == nil {
		fmt.Println("MPV process not found")
		return
	}

	if err := currentMPVProcess.Process.Kill(); err != nil {
		fmt.Printf("Failed to kill MPV process: %v\n", err)
	}
}
