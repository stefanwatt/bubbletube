package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type PlaylistModel struct {
	list             list.Model
	ID               string
	choice           string
	playbackProgress progress.Model
	volumeProgress   progress.Model
	volume           float64
	percent          float64
	duration         float64
	timePos          float64
	timeRemaining    float64
}

type model struct {
	list     list.Model
	playlist *PlaylistModel
	choice   string
	quitting bool
}

type item string

func (i item) FilterValue() string { return "" }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) SelectPlaylist() tea.Cmd {
	item, err := Find(m.list.Items(), func(i list.Item) bool {
		return i.(YTPlaylist).ID == m.choice
	})
	if err != nil {
		return nil
	}
	playlistItem := item.(YTPlaylist)
	m.playlist = &PlaylistModel{
		ID:   playlistItem.ID,
		list: MapPlaylistDetailModel(playlistItem.ID),
	}
	return nil
}

func MapPlaylistDetailModel(playlistId string) list.Model {
	items := []list.Item{}
	res := getSongsOfPlaylist(playlistId)
	for _, song := range res.Items {
		songItem := SongItem{
			ID:        song.ID,
			TitleText: song.Snippet.Title,
			VideoID:   song.Snippet.ResourceID.VideoID,
		}
		items = append(items, songItem)
	}

	l := list.New(items, songDelegate{}, defaultWidth, defaultHeight)
	l.Title = fmt.Sprintf("%d songs", len(items))
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}

func MapPlaylistModel() list.Model {
	items := []list.Item{}
	res := getPlaylists()
	for _, playlist := range res.Items {
		playlistItem := YTPlaylist{
			ID:              playlist.ID,
			TitleText:       playlist.Snippet.Title,
			DescriptionText: playlist.Snippet.Description,
		}
		items = append(items, playlistItem)
	}

	l := list.New(items, playlistDelegate{}, defaultWidth, defaultHeight)
	l.Title = "My Playlists"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}

func InitPlayingModel(m model, i SongItem) tea.Cmd {
	m.playlist.choice = string(i.ID)

	playbackProgress := progress.New(progress.WithDefaultGradient())
	playbackProgress.Full = '━'
	playbackProgress.Empty = '─'
	playbackProgress.ShowPercentage = false
	m.playlist.playbackProgress = playbackProgress

	volumeProgress := progress.New(progress.WithDefaultGradient())
	volumeProgress.Full = '━'
	volumeProgress.Empty = '─'
	cmd := volumeProgress.SetPercent(0.5)
	m.playlist.volumeProgress = volumeProgress
	return cmd
}
