package controller

import (
	view "bubbletube/view"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (sd PlaylistDelegate) Height() int                                     { return 1 }
func (sd PlaylistDelegate) Spacing() int                                    { return 0 }
func (sd PlaylistDelegate) UpdateDelegate(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (pd PlaylistDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	view.RenderPlaylist(pd, w, m, index, listItem)
}

func (sd PlaylistDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
