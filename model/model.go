package model

import (
	"fmt"
	"io"

	config "bubbletube/config"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type Delegate interface {
	Render(w io.Writer, m list.Model, index int, listItem list.Item)
	Height() int
	Spacing() int
	UpdateDelegate(_ tea.Msg, _ *list.Model) tea.Cmd
}

type PlaylistModel struct {
	List             list.Model
	PlayingSong      *SongItem
	Choice           SongItem
	ID               string
	PlaybackProgress progress.Model
	VolumeProgress   progress.Model
	Volume           float64
	Percent          float64
	Duration         float64
	TimePos          float64
	TimeRemaining    float64
	Playing          bool
}

type Screen struct {
	List         list.Model
	Playlist     *PlaylistModel
	Choice       *YTPlaylist
	Quitting     bool
	WindowWidth  int
	WindowHeight int
}

type item string

func (i item) FilterValue() string { return "" }

func (m Screen) Init() tea.Cmd {
	return nil
}

func MapPlaylistDetailModel(songDelegate list.ItemDelegate, playlistId string) list.Model {
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

	l := list.New(items, songDelegate, config.DefaultWidth, config.DefaultHeight)
	l.Title = fmt.Sprintf("%d songs", len(items))
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = config.TitleStyle
	l.Styles.PaginationStyle = config.PaginationStyle
	l.Styles.HelpStyle = config.HelpStyle
	return l
}

func MapPlaylistModel(delegate list.ItemDelegate) list.Model {
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

	l := list.New(items, delegate, config.DefaultWidth, config.DefaultHeight)
	l.Title = "My Playlists"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = config.TitleStyle
	l.Styles.PaginationStyle = config.PaginationStyle
	l.Styles.HelpStyle = config.HelpStyle
	return l
}

func InitPlayingModel(m Screen, i SongItem) tea.Cmd {
	m.Playlist.Choice = i
	m.Playlist.PlayingSong = &i

	playbackProgress := progress.New(progress.WithDefaultGradient())
	playbackProgress.Full = '━'
	playbackProgress.Empty = '─'
	playbackProgress.ShowPercentage = false
	m.Playlist.PlaybackProgress = playbackProgress

	volumeProgress := progress.New(progress.WithDefaultGradient())
	volumeProgress.Full = '━'
	volumeProgress.Empty = '─'
	cmd := volumeProgress.SetPercent(0.5)
	m.Playlist.VolumeProgress = volumeProgress
	return cmd
}
