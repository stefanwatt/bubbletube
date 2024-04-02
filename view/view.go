package view

import (
	"fmt"
	"io"
	"strings"

	config "bubbletube/config"
	model "bubbletube/model"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(config.Colors.Overlay0).
			BorderTop(true).
			BorderLeft(true).
			BorderBottom(true).
			BorderRight(true)
	centerPanelStyle = borderStyle.
				Copy().
				Padding(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(config.Colors.Peach)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
)

func View(screen model.Screen) string {
	if screen.Quitting {
		return config.QuitTextStyle.Render("Bye")
	}
	bodyStyle := lipgloss.NewStyle().
		Background(config.Colors.Base).
		Foreground(config.Colors.Text).
		Width(screen.WindowWidth).
		Height(screen.WindowHeight)

	centerPanel := ""
	width := screen.WindowWidth - 2
	songControlsHeight := 2
	songControls := borderStyle.
		Copy().
		Height(songControlsHeight).
		Width(width).
		Render(renderSongControls(screen))

	centerPanelHeight := screen.WindowHeight - songControlsHeight - 4
	switch screen.CenterPanel.(type) {
	case *model.PlaylistsPanel:
		centerPanel = centerPanel + centerPanelStyle.
			Copy().
			Width(width).
			Height(centerPanelHeight).
			Render(renderPlaylistsPanel(screen))
	case *model.PlaylistDetailPanel:
		centerPanel = centerPanel + centerPanelStyle.
			Copy().
			Width(width).
			Height(centerPanelHeight).
			Render(renderPlaylistDetailPanel(screen))
	}
	return bodyStyle.Render(
		lipgloss.JoinVertical(lipgloss.Bottom, centerPanel, songControls),
	)
}

func renderPlaylistsPanel(screen model.Screen) string {
	playlists, ok := screen.CenterPanel.(*model.PlaylistsPanel)
	if !ok {
		panic("Expected PlaylistsPanel")
	}
	return "\n" + playlists.List.View()
}

func renderPlaylistDetailPanel(screen model.Screen) string {
	playlistDetail, ok := screen.CenterPanel.(*model.PlaylistDetailPanel)
	if !ok {
		panic("Expected PlaylistDetailPanel")
	}
	return "\n" + playlistDetail.List.View()
}

func renderSongControls(screen model.Screen) string {
	progressLabel := formatProgressLabel(screen)
	volumeProgress := formatVolumeProgress(screen)
	volumesprogressLabelWidth := 3
	playPause := getPlayPauseIcon(screen)
	screen.PlaybackControls.PlaybackProgress.Width = screen.WindowWidth - len(progressLabel) - 20 - len(playPause) - volumesprogressLabelWidth
	playbackProgressLine := "\n" + playPause + screen.PlaybackControls.PlaybackProgress.View()
	songTitle := renderSongTitle(screen)
	return songTitle + playbackProgressLine + progressLabel + volumeProgress
}

func renderSongTitle(screen model.Screen) string {
	style := lipgloss.NewStyle().Foreground(config.Colors.Flamingo)
	if screen.PlaybackControls.PlayingSong == nil {
		return ""
	}
	return style.Render("  " + screen.PlaybackControls.PlayingSong.TitleText)
}

func formatProgressLabel(screen model.Screen) string {
	minutesPassed, secondsPassed := toMinutesAndSeconds(int(screen.PlaybackControls.TimePos))
	minutesDuration, secondsDuration := toMinutesAndSeconds(int(screen.PlaybackControls.Duration))

	return fmt.Sprintf(" %d:%02d/%d:%02d ", minutesPassed, secondsPassed, minutesDuration, secondsDuration)
}

func toMinutesAndSeconds(totalSeconds int) (int, int) {
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return minutes, seconds
}

func formatVolumeProgress(screen model.Screen) string {
	screen.PlaybackControls.VolumeProgress.Width = config.DefaultVolumeWidth
	return screen.PlaybackControls.VolumeProgress.View()
}

func getPlayPauseIcon(screen model.Screen) string {
	style := lipgloss.NewStyle().Width(3).Align(lipgloss.Center)
	if screen.PlaybackControls.PlayingSong == nil {
		return style.Render("")
	}
	if screen.PlaybackControls.Playing {
		return style.Render("")
	}
	return style.Render("")
}

func RenderPlaylist(d list.ItemDelegate, w io.Writer, m list.Model, index int, listItem list.Item) {
	ytItem, ok := listItem.(model.Playlist)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, ytItem.Title())

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedItemStyle.Render(" " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}

func RenderSong(d list.ItemDelegate, w io.Writer, m list.Model, index int, listItem list.Item) {
	songItem, ok := listItem.(model.SongItem)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, songItem.TitleText)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedItemStyle.Render(" " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}
