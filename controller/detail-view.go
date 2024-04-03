package controller

import (
	model "bubbletube/model"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type detailviewKeymap struct {
	back       key.Binding
	choose     key.Binding
	cursorDown key.Binding
	cursorUp   key.Binding
}

var detailviewKeys = &detailviewKeymap{
	cursorDown: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "Down"),
	),
	cursorUp: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "Up"),
	),
	back: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "Back"),
	),
	choose: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select Song"),
	),
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
		switch {
		case key.Matches(msg, detailviewKeys.cursorDown):
			detailPanel.List.CursorDown()
			selectedItem, ok := detailPanel.List.SelectedItem().(model.SongItem)
			if ok {
				detailPanel.Choice = selectedItem
			}
			return sc, nil
		case key.Matches(msg, detailviewKeys.cursorUp):
			detailPanel.List.CursorUp()
			selectedItem, ok := detailPanel.List.SelectedItem().(model.SongItem)
			if ok {
				detailPanel.Choice = selectedItem
			}
			return sc, nil

		case key.Matches(msg, detailviewKeys.back):
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
		case key.Matches(msg, detailviewKeys.choose):
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
