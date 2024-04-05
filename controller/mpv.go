package controller

import (
	model "bubbletube/model"
	"math"

	tea "github.com/charmbracelet/bubbletea"
)

func MpvFloatValueUpdate(msg model.MPVFloatValueChangedEvent, sc *ScreenController) (*ScreenController, tea.Cmd) {
	switch msg.ID {
	case model.MyMpvProperties.Volume.ID:
		percent := msg.Value / 100
		sc.Screen.PlaybackControls.Volume = math.Floor(msg.Value)
		cmd := sc.Screen.PlaybackControls.VolumeProgress.SetPercent(percent)
		return sc, cmd
	case model.MyMpvProperties.Duration.ID:
		sc.Screen.PlaybackControls.Duration = msg.Value
	case model.MyMpvProperties.PercentPos.ID:
		percent := msg.Value / 100
		sc.Screen.PlaybackControls.Percent = math.Floor(msg.Value)
		cmd := sc.Screen.PlaybackControls.PlaybackProgress.SetPercent(percent)
		return sc, cmd
	case model.MyMpvProperties.TimePos.ID:
		sc.Screen.PlaybackControls.TimePos = msg.Value
	case model.MyMpvProperties.TimeRemaining.ID:
		sc.Screen.PlaybackControls.TimeRemaining = msg.Value
		if msg.Value == 0 {
			sc.Screen.PlaybackControls = model.NextSong(sc.Screen.PlaybackControls)
		}
	}
	return sc, nil
}
