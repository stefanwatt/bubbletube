package main

import (
	"math"

	"github.com/charmbracelet/bubbles/progress"
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
		m.playlist.list.SetWidth(msg.Width)
		m.playlist.playbackProgress.Width = msg.Width
		return m, nil
	case MPVEventFloat:
		switch msg.ID {
		case 1:
			m.playlist.volume = msg.Value
		case 2:
			m.playlist.duration = msg.Value
		case 3:
			percent := msg.Value / 100
			m.playlist.percent = math.Floor(msg.Value)
			cmd := m.playlist.playbackProgress.SetPercent(percent)
			return m, cmd
		case 4:
			m.playlist.timePos = msg.Value
		case 5:
			m.playlist.timeRemaining = msg.Value
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.playlist.playbackProgress.Update(msg)
		m.playlist.playbackProgress = progressModel.(progress.Model)
		return m, cmd

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			m.choice = ""
			KillMpv()
			return m, tea.Quit
		case "ctrl+down":
			go VolumeDown()
			m.playlist.volume = m.playlist.volume - 5
			return m, nil
		case "ctrl+up":
			go VolumeUp()
			m.playlist.volume = m.playlist.volume + 5
			return m, nil
		case "left":
			go SkipBackward()
			m.playlist.timePos = m.playlist.timePos - 10
			return m, nil
		case "right":
			go SkipForward()
			m.playlist.timePos = m.playlist.timePos + 10
			return m, nil
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
			i, ok := m.playlist.list.SelectedItem().(SongItem)
			if ok {
				m.playlist.choice = string(i.ID)
				m.playlist.playbackProgress = progress.New(progress.WithDefaultGradient())
				m.playlist.playbackProgress.Full = '━'
				m.playlist.playbackProgress.Empty = '─'
				m.playlist.playbackProgress.ShowPercentage = false
				SelectSong(i)
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.playlist.list.Update(msg)
	m.choice = ""
	return m, cmd
}
