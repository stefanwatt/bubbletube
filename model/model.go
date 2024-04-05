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

type CenterPanel interface {
	GetChoice() list.Item
	SetChoice(list.Item)
	BubbletubeCenterPanel()
	GetList() *list.Model
	SetList(list.Model)
}

type PlaylistDetailPanel struct {
	List   list.Model
	Choice SongItem
	ID     string
}

func (p PlaylistDetailPanel) GetList() *list.Model  { return &p.List }
func (p *PlaylistDetailPanel) SetList(l list.Model) { p.List = l }

func (p PlaylistDetailPanel) GetChoice() list.Item { return p.Choice }
func (p *PlaylistDetailPanel) SetChoice(i list.Item) {
	if songItem, ok := i.(SongItem); ok {
		p.Choice = songItem
	}
}
func (PlaylistDetailPanel) BubbletubeCenterPanel() {}

type PlaylistsPanel struct {
	List   list.Model
	Choice Playlist
}

func (p PlaylistsPanel) GetList() *list.Model  { return &p.List }
func (p *PlaylistsPanel) SetList(l list.Model) { p.List = l }
func (p PlaylistsPanel) GetChoice() list.Item  { return p.Choice }
func (p *PlaylistsPanel) SetChoice(i list.Item) {
	if playlists, ok := i.(Playlist); ok {
		p.Choice = playlists
	}
}
func (PlaylistsPanel) BubbletubeCenterPanel() {}

type PlaybackControls struct {
	PlayingSong      *SongItem
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
	QueuePanel       QueuePanel
	CenterPanel      CenterPanel
	PlaybackControls PlaybackControls
	Quitting         bool
	WindowWidth      int
	WindowHeight     int
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
	l.Styles.Title = config.TitleStyle
	l.Styles.PaginationStyle = config.PaginationStyle
	l.Styles.HelpStyle = config.HelpStyle
	return l
}

func MapPlaylistsModel(delegate list.ItemDelegate) list.Model {
	items := []list.Item{}
	res := getPlaylists()
	for _, playlist := range res.Items {
		playlistItem := Playlist{
			ID:              playlist.ID,
			TitleText:       playlist.Snippet.Title,
			DescriptionText: playlist.Snippet.Description,
		}
		items = append(items, playlistItem)
	}

	l := list.New(items, delegate, config.DefaultWidth, config.DefaultHeight)
	l.Title = "My Playlists"
	l.SetShowStatusBar(false)
	l.Styles.Title = config.TitleStyle
	l.Styles.PaginationStyle = config.PaginationStyle
	l.Styles.HelpStyle = config.HelpStyle
	return l
}

func MapDefaultPlaybackProgress() progress.Model {
	playbackProgress := progress.New(progress.WithDefaultGradient())
	playbackProgress.Full = '━'
	playbackProgress.Empty = '─'
	playbackProgress.ShowPercentage = false
	return playbackProgress
}

func MapDefaultVolumeProgress() progress.Model {
	volumeProgress := progress.New(progress.WithDefaultGradient())
	volumeProgress.Full = '━'
	volumeProgress.Empty = '─'
	return volumeProgress
}

func InitPlayingModel(screen *Screen, detailPanel *PlaylistDetailPanel, i SongItem) tea.Cmd {
	screen.CenterPanel = &PlaylistDetailPanel{
		ID:     detailPanel.ID,
		List:   detailPanel.List,
		Choice: i,
	}

	screen.PlaybackControls.PlayingSong = &i
	screen.PlaybackControls.PlaybackProgress = MapDefaultPlaybackProgress()

	screen.PlaybackControls.VolumeProgress = MapDefaultVolumeProgress()
	return screen.PlaybackControls.VolumeProgress.SetPercent(config.DefaultVolume)
}

func MapWaitlistModel(items []list.Item, songDelegate list.ItemDelegate) list.Model {
	waitlist := list.New(items, songDelegate, config.DefaultWidth, config.DefaultWaitlistHeight)
	waitlist.SetShowStatusBar(false)
	waitlist.SetShowTitle(true)
	waitlist.SetShowHelp(false)
	return waitlist
}

func MapDefaultScreen(playlists list.Model, songDelegate list.ItemDelegate) Screen {
	waitlist := MapWaitlistModel([]list.Item{}, songDelegate)
	waitlist.Title = "Waitlist"
	playlist := MapWaitlistModel([]list.Item{}, songDelegate)
	playlist.Title = "Playlist"
	return Screen{
		CenterPanel: &PlaylistsPanel{
			List: playlists,
		},
		QueuePanel: QueuePanel{
			Waitlist: waitlist,
			Playlist: playlist,
		},
		PlaybackControls: PlaybackControls{
			TimePos:          0,
			Volume:           config.DefaultVolume,
			Playing:          false,
			Duration:         0.0,
			TimeRemaining:    0.0,
			Percent:          0.0,
			VolumeProgress:   MapDefaultVolumeProgress(),
			PlayingSong:      nil,
			PlaybackProgress: MapDefaultPlaybackProgress(),
		},
		Quitting:     false,
		WindowWidth:  config.DefaultHeight,
		WindowHeight: config.DefaultWidth,
	}
}
