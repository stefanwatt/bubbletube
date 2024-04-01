package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) View() string {
	if m.playlist != nil {
		res := "\n" + m.playlist.list.View()
		if m.playlist != nil {
			totalSeconds := int(m.playlist.duration)
			minutesDuration := totalSeconds / 60
			secondsDuration := totalSeconds % 60
			totalSeconds = int(m.playlist.timePos)
			minutesPassed := totalSeconds / 60
			secondsPassed := totalSeconds % 60

			progressLabel := fmt.Sprintf(" %d:%02d/%d:%02d ",
				minutesPassed,
				secondsPassed,
				minutesDuration,
				secondsDuration,
			)

			volumeLength := 20
			m.playlist.volumeProgress.Width = volumeLength
			volumeProgress := m.playlist.volumeProgress.View()
			m.playlist.playbackProgress.Width = m.playlist.list.Width() - len(progressLabel) - volumeLength
			res = res + "\n" + m.playlist.playbackProgress.View() + progressLabel + volumeProgress + "\n"
		}
		return res
	}
	if m.quitting {
		return quitTextStyle.Render("Bye")
	}
	return "\n" + m.list.View()
}

type playlistDelegate struct{}

func (d playlistDelegate) Height() int                             { return 1 }
func (d playlistDelegate) Spacing() int                            { return 0 }
func (d playlistDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d playlistDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	ytItem, ok := listItem.(YTPlaylist)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, ytItem.Title())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}

type songDelegate struct{}

func (d songDelegate) Height() int                             { return 1 }
func (d songDelegate) Spacing() int                            { return 0 }
func (d songDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d songDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	songItem, ok := listItem.(SongItem)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, songItem.Title())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}
