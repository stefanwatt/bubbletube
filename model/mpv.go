package model

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/DexterLB/mpvipc"
	tea "github.com/charmbracelet/bubbletea"
)

type MPVFloatValueChangedEvent struct {
	ID    int64
	Value float64
}

type MPVVoidValueChangedEvent struct {
	Name   string
	Reason string
}

func mpvFloatEventCmd(event mpvipc.Event) tea.Cmd {
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

func mpvVoidEventCmd(event mpvipc.Event) tea.Cmd {
	return func() tea.Msg {
		return MPVVoidValueChangedEvent{
			Name:   event.Name,
			Reason: event.Reason,
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

	events, stopListening := conn.NewEventListener()
	go func() {
		for event := range events {
			if program == nil {
				continue
			}
			if event.Name == "end-file" {
				msg := mpvVoidEventCmd(*event)().(MPVVoidValueChangedEvent)
				program.Send(msg)
			} else {
				msg := mpvFloatEventCmd(*event)().(MPVFloatValueChangedEvent)
				switch msg.ID {
				case MyMpvProperties.TimeRemaining.ID:
					time_remaining = msg.Value
				case MyMpvProperties.TimePos.ID:
					time_pos = msg.Value
				}
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

type MpvProperty struct {
	Name string
	ID   int64
}

// Make sure all fields start with uppercase letters to export them
type MpvProperties struct {
	Volume        MpvProperty
	Duration      MpvProperty
	PercentPos    MpvProperty
	TimePos       MpvProperty
	TimeRemaining MpvProperty
}

// Example variable with exported fields
var MyMpvProperties = MpvProperties{
	Volume:        MpvProperty{"volume", 1},
	Duration:      MpvProperty{"duration", 2},
	PercentPos:    MpvProperty{"percent-pos", 3},
	TimePos:       MpvProperty{"time-pos", 4},
	TimeRemaining: MpvProperty{"time-remaining", 5},
}

func observeProperties() {
	val := reflect.ValueOf(MyMpvProperties)
	for i := 0; i < val.NumField(); i++ {
		property := val.Field(i).Interface().(MpvProperty)
		if _, err := conn.Call("observe_property", property.ID, property.Name); err != nil {
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
