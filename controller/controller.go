package controller

import (
	"math"

	model "bubbletube/model"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func updateListView(msg tea.Msg, sc *ScreenController) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sc.Screen.WindowWidth = msg.Width
		sc.Screen.WindowHeight = msg.Height
		sc.Screen.List.SetWidth(80)
		return sc, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			sc.Screen.Quitting = true
			model.KillMpv()
			return sc, tea.Quit

		case "enter":
			i, ok := sc.Screen.List.SelectedItem().(model.YTPlaylist)
			if ok {
				sc.Screen.Choice = &i
			}
			l := model.MapPlaylistDetailModel(sc.SongDelegate, sc.Screen.Choice.ID)
			sc.Screen.Playlist = &model.PlaylistModel{
				List: l,
				ID:   sc.Screen.Choice.ID,
			}
			return sc, nil
		}
	}

	var cmd tea.Cmd
	sc.Screen.List, cmd = sc.Screen.List.Update(msg)
	sc.Screen.Choice = nil
	return sc, cmd
}

func updateDetailView(msg tea.Msg, sc *ScreenController) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sc.Screen.WindowWidth = msg.Width
		sc.Screen.WindowHeight = msg.Height
		sc.Screen.Playlist.List.SetWidth(80)
		sc.Screen.Playlist.PlaybackProgress.Width = 80
		return sc, nil
	case model.MPVEventFloat:
		switch msg.ID {
		case 1:
			percent := msg.Value / 100
			sc.Screen.Playlist.Volume = math.Floor(msg.Value)
			cmd := sc.Screen.Playlist.VolumeProgress.SetPercent(percent)
			return sc, cmd
		case 2:
			sc.Screen.Playlist.Duration = msg.Value
		case 3:
			percent := msg.Value / 100
			sc.Screen.Playlist.Percent = math.Floor(msg.Value)
			cmd := sc.Screen.Playlist.PlaybackProgress.SetPercent(percent)
			return sc, cmd
		case 4:
			sc.Screen.Playlist.TimePos = msg.Value
		case 5:
			sc.Screen.Playlist.TimeRemaining = msg.Value
		}
		return sc, nil

	case progress.FrameMsg:
		var (
			cmds         []tea.Cmd
			cmd          tea.Cmd
			updatedModel tea.Model
		)

		updatedModel, cmd = sc.Screen.Playlist.PlaybackProgress.Update(msg)
		cmds = append(cmds, cmd)
		sc.Screen.Playlist.PlaybackProgress = updatedModel.(progress.Model)
		updatedModel, cmd = sc.Screen.Playlist.VolumeProgress.Update(msg)
		sc.Screen.Playlist.VolumeProgress = updatedModel.(progress.Model)
		cmds = append(cmds, cmd)
		return sc, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			sc.Screen.Quitting = true
			sc.Screen.Choice = nil
			model.KillMpv()
			return sc, tea.Quit
		case "ctrl+down":
			model.VolumeDown()
			sc.Screen.Playlist.Volume = sc.Screen.Playlist.Volume - 5
			return sc, nil
		case "ctrl+up":
			model.VolumeUp()
			sc.Screen.Playlist.Volume = sc.Screen.Playlist.Volume + 5
			return sc, nil
		case "left":
			model.SkipBackward()
			sc.Screen.Playlist.TimePos = sc.Screen.Playlist.TimePos - 10
			return sc, nil
		case "right":
			model.SkipForward()
			sc.Screen.Playlist.TimePos = sc.Screen.Playlist.TimePos + 10
			return sc, nil
		case "down":
			sc.Screen.Playlist.List.CursorDown()
			selectedItem, ok := sc.Screen.Playlist.List.SelectedItem().(model.SongItem)
			if ok {
				sc.Screen.Playlist.Choice = selectedItem
			}
			return sc, nil
		case "up":
			sc.Screen.Playlist.List.CursorUp()
			selectedItem, ok := sc.Screen.Playlist.List.SelectedItem().(model.SongItem)
			if ok {
				sc.Screen.Playlist.Choice = selectedItem
			}
			return sc, nil
		case "p":
			sc.Screen.Playlist.Playing = !model.TogglePlayback()
			return sc, nil
		case "backspace":
			sc.Screen.Playlist = nil
			return sc, nil
		case "enter":
			i, ok := sc.Screen.Playlist.List.SelectedItem().(model.SongItem)
			var cmd tea.Cmd
			if ok {
				cmd = model.InitPlayingModel(sc.Screen, i)
				model.SelectSong(i)
			}
			return sc, cmd
		}
	}

	var cmd tea.Cmd
	sc.Screen.List, cmd = sc.Screen.Playlist.List.Update(msg)
	sc.Screen.Choice = nil
	return sc, cmd
}
