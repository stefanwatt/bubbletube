package main

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

type MPVEventFloat struct {
	ID    int64
	Value float64
}

func mpvEventCmd(event mpvipc.Event) tea.Cmd {
	return func() tea.Msg {
		if event.Data == nil {
			return MPVEventFloat{
				ID:    event.ID,
				Value: 0,
			}
		}
		return MPVEventFloat{
			ID:    event.ID,
			Value: event.Data.(float64),
		}
	}
}

var (
	conn           *mpvipc.Connection
	playing        = false
	volume         = 0.0
	time_remaining = 0.0
	time_pos       = 0.0
)

func initConn() {
	conn = mpvipc.NewConnection("/tmp/mpv_socket")
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}
	events, stopListening := conn.NewEventListener()
	_, err = conn.Call("observe_property", 1, "volume")
	if err != nil {
		fmt.Print(err)
	}

	_, err = conn.Call("observe_property", 2, "duration")
	if err != nil {
		fmt.Print(err)
	}
	_, err = conn.Call("observe_property", 3, "percent-pos")
	if err != nil {
		fmt.Print(err)
	}
	_, err = conn.Call("observe_property", 4, "time-pos")
	if err != nil {
		fmt.Print(err)
	}
	_, err = conn.Call("observe_property", 5, "time-remaining")
	if err != nil {
		fmt.Print(err)
	}
	go func() {
		conn.WaitUntilClosed()
		stopListening <- struct{}{}
	}()

	for event := range events {
		switch event.ID {
		case 1, 2, 3, 4, 5:
			msg := mpvEventCmd(*event)().(MPVEventFloat)
			program.Send(msg)
		}

		if event.Data == nil {
			return
		}
		switch event.ID {
		case 1:
			volume = event.Data.(float64)
		case 4:
			time_pos = event.Data.(float64)
		case 5:
			time_remaining = event.Data.(float64)
		}
	}
}

func TogglePlayback() {
	playing = !playing
	err := conn.Set("pause", playing)
	if err != nil {
		log.Fatal(err)
	}
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

var (
	currentMPVProcess *exec.Cmd
	started           = false
)

func KillMpv() {
	if currentMPVProcess == nil || currentMPVProcess.Process == nil {
		fmt.Println("MPV process not found")
		return
	}

	if err := currentMPVProcess.Process.Kill(); err != nil {
		fmt.Printf("Failed to kill MPV process: %v\n", err)
	}
}

func PlaySong(item SongItem) {
	started = true
	currentMPVProcess = exec.Command(
		"mpv",
		"--input-unix-socket=/tmp/mpv_socket",
		"--no-video",
		"https://www.youtube.com/watch?v="+item.VideoID,
	)
	if err := currentMPVProcess.Start(); err != nil {
		fmt.Printf("Failed to start MPV: %v\n", err)
		return
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			sig := <-sigChan
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				KillMpv()
			}
		}
	}()

	if conn == nil {
		time.Sleep(1 * time.Second)
		go func() {
			initConn()
		}()
	}
	go func() {
		if err := currentMPVProcess.Wait(); err != nil {
			fmt.Printf("MPV process exited with error: %v\n", err)
		} else {
			fmt.Println("MPV process exited normally")
		}
	}()
}
