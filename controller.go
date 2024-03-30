package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.playlist != nil {
		return updateDetailView(msg, m)
	} else {
		return updateListView(msg, m)
	}
}

func updateListView(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			KillMpv()
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(YTPlaylist)
			if ok {
				m.choice = string(i.ID)
			}
			l := MapPlaylistDetailModel(m.choice)
			m.playlist = &PlaylistModel{
				list: l,
				ID:   m.choice,
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	m.choice = ""
	return m, cmd
}

func updateDetailView(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			m.choice = ""
			KillMpv()
			return m, tea.Quit
		case "down":
			m.playlist.list.CursorDown()
			m.choice = ""
			return m, nil
		case "up":
			m.playlist.list.CursorUp()
			m.choice = ""
			return m, nil
		case "p":
			TogglePlayback()
			return m, nil
		case "backspace":
			m.playlist = nil
			m.choice = ""
			return m, nil
		case "enter":
			i, ok := m.playlist.list.SelectedItem().(SongItem) // TODO this is the wrong model
			if ok {
				m.playlist.choice = string(i.ID)
				PlaySong(i)
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.playlist.list.Update(msg)
	m.choice = ""
	return m, cmd
}
