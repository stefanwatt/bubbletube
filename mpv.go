package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/DexterLB/mpvipc"
)

var (
	conn    *mpvipc.Connection
	playing = false
)

func initConn() {
	conn = mpvipc.NewConnection("/tmp/mpv_socket")
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Call("observe_property", 42, "volume")
	if err != nil {
		fmt.Print(err)
	}
}

func TogglePlayback() {
	err := conn.Set("pause", playing)
	playing = !playing
	if err != nil {
		log.Fatal(err)
	}
}

var (
	currentMPVProcess *exec.Cmd
	started           = false
)

func PlaySong(item PlaylistItem) {
	if started || (currentMPVProcess != nil && currentMPVProcess.Process != nil) {
		fmt.Println("Already playing")
		return
	}
	fmt.Println("Starting playback")
	started = true
	currentMPVProcess := exec.Command("mpv", "--input-unix-socket=/tmp/mpv_socket", "--no-video", "https://www.youtube.com/watch?v="+item.Snippet.ResourceID.VideoID)
	currentMPVProcess.Start()
	if conn == nil {
		time.Sleep(1 * time.Second)
		initConn()
	}
}
