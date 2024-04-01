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
		sc.Screen.CenterPanel.GetList().SetWidth(80)
		return sc, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			sc.Screen.Quitting = true
			model.KillMpv()
			return sc, tea.Quit

		case "enter":
			i, ok := sc.Screen.CenterPanel.GetList().SelectedItem().(model.Playlist)
			if ok {
				sc.Screen.CenterPanel.SetChoice(i)
			}
			var selectedPlaylist model.Playlist
			selectedPlaylist, ok = sc.Screen.CenterPanel.GetChoice().(model.Playlist)
			if !ok {
				panic("Failed to cast to SongItem")
			}
			l := model.MapPlaylistDetailModel(sc.SongDelegate, selectedPlaylist.ID)
			choice, ok := l.Items()[0].(model.SongItem)
			if !ok {
				panic("Failed to cast to SongItem")
			}
			sc.Screen.CenterPanel = &model.PlaylistDetailPanel{
				List:   l,
				Choice: choice,
				ID:     selectedPlaylist.ID,
			}
			return sc, nil
		}
	}

	playlistsPanel, ok := sc.Screen.CenterPanel.(*model.PlaylistsPanel)
	if !ok {
		panic("Failed to cast to PlaylistsPanel")
	}
	updatedList, cmd := playlistsPanel.List.Update(msg)
	sc.Screen.CenterPanel.SetList(updatedList)
	sc.Screen.CenterPanel.SetChoice(nil)
	return sc, cmd
}

func updateDetailView(msg tea.Msg, sc *ScreenController) (tea.Model, tea.Cmd) {
	detailPanel, ok := sc.Screen.CenterPanel.(*model.PlaylistDetailPanel)
	if !ok {
		panic("Failed to cast to PlaylistDetailPanel")
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sc.Screen.WindowWidth = msg.Width
		sc.Screen.WindowHeight = msg.Height
		detailPanel.List.SetWidth(80)
		sc.Screen.PlaybackControls.PlaybackProgress.Width = 80
		return sc, nil
	case model.MPVEventFloat:
		switch msg.ID {
		case 1:
			percent := msg.Value / 100
			sc.Screen.PlaybackControls.Volume = math.Floor(msg.Value)
			cmd := sc.Screen.PlaybackControls.VolumeProgress.SetPercent(percent)
			return sc, cmd
		case 2:
			sc.Screen.PlaybackControls.Duration = msg.Value
		case 3:
			percent := msg.Value / 100
			sc.Screen.PlaybackControls.Percent = math.Floor(msg.Value)
			cmd := sc.Screen.PlaybackControls.PlaybackProgress.SetPercent(percent)
			return sc, cmd
		case 4:
			sc.Screen.PlaybackControls.TimePos = msg.Value
		case 5:
			sc.Screen.PlaybackControls.TimeRemaining = msg.Value
		}
		return sc, nil

	case progress.FrameMsg:
		var (
			cmds         []tea.Cmd
			cmd          tea.Cmd
			updatedModel tea.Model
		)

		updatedModel, cmd = sc.Screen.PlaybackControls.PlaybackProgress.Update(msg)
		cmds = append(cmds, cmd)
		sc.Screen.PlaybackControls.PlaybackProgress = updatedModel.(progress.Model)
		updatedModel, cmd = sc.Screen.PlaybackControls.VolumeProgress.Update(msg)
		sc.Screen.PlaybackControls.VolumeProgress = updatedModel.(progress.Model)
		cmds = append(cmds, cmd)
		return sc, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			sc.Screen.Quitting = true
			detailPanel.SetChoice(nil)
			model.KillMpv()
			return sc, tea.Quit
		case "ctrl+down":
			model.VolumeDown()
			sc.Screen.PlaybackControls.Volume = sc.Screen.PlaybackControls.Volume - 5
			return sc, nil
		case "ctrl+up":
			model.VolumeUp()
			sc.Screen.PlaybackControls.Volume = sc.Screen.PlaybackControls.Volume + 5
			return sc, nil
		case "left":
			model.SkipBackward()
			sc.Screen.PlaybackControls.TimePos = sc.Screen.PlaybackControls.TimePos - 10
			return sc, nil
		case "right":
			model.SkipForward()
			sc.Screen.PlaybackControls.TimePos = sc.Screen.PlaybackControls.TimePos + 10
			return sc, nil
		case "down":
			detailPanel.List.CursorDown()
			selectedItem, ok := detailPanel.List.SelectedItem().(model.SongItem)
			if ok {
				detailPanel.Choice = selectedItem
			}
			return sc, nil
		case "up":
			detailPanel.List.CursorUp()
			selectedItem, ok := detailPanel.List.SelectedItem().(model.SongItem)
			if ok {
				detailPanel.Choice = selectedItem
			}
			return sc, nil
		case "p":
			sc.Screen.PlaybackControls.Playing = !model.TogglePlayback()
			return sc, nil
		case "backspace":
			list := model.MapPlaylistsModel(sc.PlaylistDelegate)
			choice, ok := list.Items()[0].(model.Playlist)
			if !ok {
				panic("Failed to cast to Playlist")
			}
			sc.Screen.CenterPanel = &model.PlaylistsPanel{
				List:   list,
				Choice: choice,
			}
			return sc, nil
		case "enter":
			i, ok := detailPanel.List.SelectedItem().(model.SongItem)
			var cmd tea.Cmd
			if ok {
				cmd = model.InitPlayingModel(&sc.Screen, detailPanel, i)
				model.SelectSong(i)
			}
			return sc, cmd
		}
	}

	playlistDetailPanel, ok := sc.Screen.CenterPanel.(*model.PlaylistDetailPanel)
	if !ok {
		panic("Failed to cast to PlaylistsPanel")
	}
	updatedList, cmd := playlistDetailPanel.List.Update(msg)
	sc.Screen.CenterPanel.SetList(updatedList)
	sc.Screen.CenterPanel.SetChoice(nil)
	return sc, cmd
}
