package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
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
	playing = !playing
	err := conn.Set("pause", playing)
	if err != nil {
		log.Fatal(err)
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
		initConn()
	}
	go func() {
		if err := currentMPVProcess.Wait(); err != nil {
			fmt.Printf("MPV process exited with error: %v\n", err)
		} else {
			fmt.Println("MPV process exited normally")
		}
	}()
}
