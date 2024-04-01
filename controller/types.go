package controller

import (
	model "bubbletube/model"
	view "bubbletube/view"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ScreenUpdater interface {
	Update(tea.Msg) (tea.Model, tea.Cmd)
}

type ScreenController struct {
	PlaylistDelegate PlaylistDelegate
	SongDelegate     SongDelegate
	Screen           model.Screen
}

// Define separate structs for each delegate type.
type PlaylistDelegate struct {
	*ScreenController
}

type SongDelegate struct {
	*ScreenController
}

func NewScreenController(screen model.Screen) *ScreenController {
	sc := &ScreenController{Screen: screen}
	// Initialize delegates with a reference back to the ScreenController.
	sc.PlaylistDelegate = PlaylistDelegate{sc}
	sc.SongDelegate = SongDelegate{sc}
	return sc
}

func (sc *ScreenController) View() string {
	return view.View(sc.Screen)
}

func (sc *ScreenController) Init() tea.Cmd {
	return nil
}

func (sc *ScreenController) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sc.Screen.Playlist != nil {
		return updateDetailView(msg, sc)
	}
	return updateListView(msg, sc)
}

func (sd PlaylistDelegate) Height() int                                     { return 1 }
func (sd PlaylistDelegate) Spacing() int                                    { return 0 }
func (sd PlaylistDelegate) UpdateDelegate(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (pd PlaylistDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	view.RenderPlaylist(pd, w, m, index, listItem)
}

func (sd PlaylistDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (sd SongDelegate) Height() int                                     { return 1 }
func (sd SongDelegate) Spacing() int                                    { return 0 }
func (sd SongDelegate) UpdateDelegate(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (sd SongDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	view.RenderSong(sd, w, m, index, listItem)
}

func (sd SongDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
