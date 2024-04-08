package controller

import (
	model "bubbletube/model"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type listviewKeymap struct {
	choose key.Binding
	noop   key.Binding
}

var listviewKeys = &listviewKeymap{
	choose: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select Playlist"),
	),
	noop: key.NewBinding(
		key.WithKeys("q"),
	),
}

func updateListView(msg tea.Msg, sc *ScreenController) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if sc.Screen.CenterPanel.GetList().FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, listviewKeys.noop):
			return sc, nil
		case key.Matches(msg, listviewKeys.choose):
			currentList := sc.Screen.CenterPanel.GetList()
			i, ok := currentList.SelectedItem().(model.Playlist)
			if ok {
				sc.Screen.CenterPanel.SetChoice(i)
				var selectedPlaylist model.Playlist
				selectedPlaylist, ok = sc.Screen.CenterPanel.GetChoice().(model.Playlist)
				if !ok {
					panic("Failed to cast to SongItem")
				}
				l := model.MapPlaylistDetailModel(sc.SongDelegate, selectedPlaylist)
				choice, ok := l.Items()[0].(model.SongItem)
				l.SetHeight(currentList.Height())
				if !ok {
					panic("Failed to cast to SongItem")
				}
				sc.Screen.CenterPanel = &model.PlaylistDetailPanel{
					List:   l,
					Choice: choice,
					ID:     selectedPlaylist.ID,
				}
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
