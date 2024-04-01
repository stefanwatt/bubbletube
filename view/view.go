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

func View(m model.Screen) string {
	padding := getPaddingLeft(m)
	if m.Playlist != nil {
		return applyPaddingToView("\n"+renderPlaylist(m), padding)
	}
	if m.Quitting {
		return config.QuitTextStyle.Render("Bye")
	}
	return applyPaddingToView("\n"+m.List.View(), padding)
}

func applyPaddingToView(view string, paddingLeft string) string {
	lines := strings.Split(view, "\n") // Split the view into lines
	for i, line := range lines {
		lines[i] = paddingLeft + line // Prepend padding to each line
	}
	return strings.Join(lines, "\n") // Join the lines back together
}

func getPaddingLeft(m model.Screen) string {
	spaces := 0
	if m.WindowWidth > 80 {
		spaces = int(math.Floor(float64(m.WindowWidth-80) / 2))
	}
	return strings.Repeat(" ", spaces)
}

func renderPlaylist(m model.Screen) string {
	res := "\n" + m.Playlist.List.View()
	if m.Playlist.PlayingSong != nil {
		res += renderSongControls(m)
	}
	return res
}

func renderSongControls(m model.Screen) string {
	progressLabel := formatProgressLabel(m)
	volumeProgress := formatVolumeProgress(m)
	playPause := getPlayPauseIcon(m)

	m.Playlist.PlaybackProgress.Width = m.Playlist.List.Width() - len(progressLabel) - 20 - len(playPause)
	playbackProgressLine := "\n" + playPause + m.Playlist.PlaybackProgress.View()

	return "\n  " + m.Playlist.PlayingSong.TitleText + playbackProgressLine + progressLabel + volumeProgress
}

func formatProgressLabel(m model.Screen) string {
	minutesPassed, secondsPassed := formatTime(int(m.Playlist.TimePos))
	minutesDuration, secondsDuration := formatTime(int(m.Playlist.Duration))

	return fmt.Sprintf(" %d:%02d/%d:%02d ", minutesPassed, secondsPassed, minutesDuration, secondsDuration)
}

func formatTime(totalSeconds int) (int, int) {
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return minutes, seconds
}

func formatVolumeProgress(m model.Screen) string {
	m.Playlist.VolumeProgress.Width = 20 // Consider making this a constant or configurable field
	return m.Playlist.VolumeProgress.View()
}

func getPlayPauseIcon(m model.Screen) string {
	if m.Playlist.PlayingSong == nil {
		return ""
	}
	if m.Playlist.Playing {
		return "  "
	}
	return "  "
}

func RenderPlaylist(d list.ItemDelegate, w io.Writer, m list.Model, index int, listItem list.Item) {
	ytItem, ok := listItem.(model.YTPlaylist)
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
