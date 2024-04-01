package controller

import (
	model "bubbletube/model"

	"github.com/charmbracelet/bubbles/list"
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
		if sc.Screen.CenterPanel.GetList().FilterState() == list.Filtering {
			break
		}
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

	case tea.KeyMsg:
		if detailPanel.List.FilterState() == list.Filtering {
			break
		}
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			sc.Screen.Quitting = true
			detailPanel.SetChoice(nil)
			model.KillMpv()
			return sc, tea.Quit
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
