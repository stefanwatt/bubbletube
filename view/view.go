package view

import (
	"fmt"
	"io"
	"math"
	"strings"

	config "bubbletube/config"
	model "bubbletube/model"

	"github.com/charmbracelet/bubbles/list"
)

func View(screen model.Screen) string {
	if screen.Quitting {
		return config.QuitTextStyle.Render("Bye")
	}
	padding := getPaddingLeft(screen)
	res := ""
	switch screen.CenterPanel.(type) {
	case *model.PlaylistsPanel:
		res = res + applyPaddingToView(renderPlaylistsPanel(screen), padding)
	case *model.PlaylistDetailPanel:
		res = res + applyPaddingToView(renderPlaylistDetailPanel(screen), padding)
	}
	return res + applyPaddingToView(renderSongControls(screen), padding)
}

func applyPaddingToView(view string, paddingLeft string) string {
	lines := strings.Split(view, "\n") // Split the view into lines
	for i, line := range lines {
		lines[i] = paddingLeft + line // Prepend padding to each line
	}
	return strings.Join(lines, "\n") // Join the lines back together
}

func getPaddingLeft(screen model.Screen) string {
	spaces := 0
	if screen.WindowWidth > 80 {
		spaces = int(math.Floor(float64(screen.WindowWidth-80) / 2))
	}
	return strings.Repeat(" ", spaces)
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
	playPause := getPlayPauseIcon(screen)
	screen.PlaybackControls.PlaybackProgress.Width = screen.CenterPanel.GetList().Width() - len(progressLabel) - 20 - len(playPause)
	playbackProgressLine := "\n" + playPause + screen.PlaybackControls.PlaybackProgress.View()
	songTitle := renderSongTitle(screen)
	return "\n" + songTitle + playbackProgressLine + progressLabel + volumeProgress
}

func renderSongTitle(screen model.Screen) string {
	if screen.PlaybackControls.PlayingSong == nil {
		return ""
	}
	return "  " + screen.PlaybackControls.PlayingSong.TitleText
}

func formatProgressLabel(screen model.Screen) string {
	minutesPassed, secondsPassed := formatTime(int(screen.PlaybackControls.TimePos))
	minutesDuration, secondsDuration := formatTime(int(screen.PlaybackControls.Duration))

	return fmt.Sprintf(" %d:%02d/%d:%02d ", minutesPassed, secondsPassed, minutesDuration, secondsDuration)
}

func formatTime(totalSeconds int) (int, int) {
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return minutes, seconds
}

func formatVolumeProgress(screen model.Screen) string {
	screen.PlaybackControls.VolumeProgress.Width = 20 // Consider making this a constant or configurable field
	return screen.PlaybackControls.VolumeProgress.View()
}

func getPlayPauseIcon(screen model.Screen) string {
	if screen.PlaybackControls.PlayingSong == nil {
		return ""
	}
	if screen.PlaybackControls.Playing {
		return "  "
	}
	return "  "
}

func RenderPlaylist(d list.ItemDelegate, w io.Writer, m list.Model, index int, listItem list.Item) {
	ytItem, ok := listItem.(model.Playlist)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, ytItem.Title())

	fn := config.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return config.SelectedItemStyle.Render("> " + strings.Join(s, " "))
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

	fn := config.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return config.SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}
