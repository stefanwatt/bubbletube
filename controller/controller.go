package controller

import (
	model "bubbletube/model"
	view "bubbletube/view"
	"io"
	"math"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+down":
			model.VolumeDown()
			sc.Screen.PlaybackControls.Volume = sc.Screen.PlaybackControls.Volume - 5
			return sc, nil
		case "ctrl+up":
			model.VolumeUp()
			sc.Screen.PlaybackControls.Volume = sc.Screen.PlaybackControls.Volume + 5
			return sc, nil
		case "ctrl+right":
			sc.Screen.PlaybackControls = model.NextSong(sc.Screen.PlaybackControls)
			return sc, nil
		case "left":
			model.SkipBackward()
			sc.Screen.PlaybackControls.TimePos = sc.Screen.PlaybackControls.TimePos - 10
			return sc, nil
		case "right":
			model.SkipForward()
			sc.Screen.PlaybackControls.TimePos = sc.Screen.PlaybackControls.TimePos + 10
			return sc, nil
		case "p":
			sc.Screen.PlaybackControls.Playing = !model.TogglePlayback()
			return sc, nil
		}

	case model.MPVFloatValueChangedEvent:
		switch msg.ID {
		case 1:
			percent := msg.Value / 100
			sc.Screen.PlaybackControls.Volume = math.Floor(msg.Value)
			cmd := sc.Screen.PlaybackControls.VolumeProgress.SetPercent(percent)
			return sc, cmd
		case 2:
			sc.Screen.PlaybackControls.Duration = msg.Value
		case 3:
			percent := msg.Value / 100
			sc.Screen.PlaybackControls.Percent = math.Floor(msg.Value)
			cmd := sc.Screen.PlaybackControls.PlaybackProgress.SetPercent(percent)
			return sc, cmd
		case 4:
			sc.Screen.PlaybackControls.TimePos = msg.Value
		case 5:
			sc.Screen.PlaybackControls.TimeRemaining = msg.Value
			if msg.Value == 0 {
				sc.Screen.PlaybackControls = model.NextSong(sc.Screen.PlaybackControls)
			}
		}
		return sc, nil

	case progress.FrameMsg:
		var (
			cmds         []tea.Cmd
			cmd          tea.Cmd
			updatedModel tea.Model
		)

		updatedModel, cmd = sc.Screen.PlaybackControls.PlaybackProgress.Update(msg)
		cmds = append(cmds, cmd)
		sc.Screen.PlaybackControls.PlaybackProgress = updatedModel.(progress.Model)
		updatedModel, cmd = sc.Screen.PlaybackControls.VolumeProgress.Update(msg)
		sc.Screen.PlaybackControls.VolumeProgress = updatedModel.(progress.Model)
		cmds = append(cmds, cmd)
		return sc, tea.Batch(cmds...)
	}

	switch sc.Screen.CenterPanel.(type) {
	case *model.PlaylistDetailPanel:
		return updateDetailView(msg, sc)
	case *model.PlaylistsPanel:
		return updateListView(msg, sc)
	}
	return sc, nil
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

func (sd SongDelegate) Height() int  { return 1 }
func (sd SongDelegate) Spacing() int { return 0 }
func (sd SongDelegate) ShortHelpFunc() []key.Binding {
	return []key.Binding{}
}

func (sd SongDelegate) FullHelpFunc() [][]key.Binding {
	return [][]key.Binding{}
}

type delegateKeyMap struct {
	enqueue key.Binding
}

var delegateKeys = &delegateKeyMap{
	enqueue: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Add to Queue"),
	),
}

var statusMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
	Render

func (sd SongDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	item, ok := m.SelectedItem().(model.SongItem)
	if !ok {
		return nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, delegateKeys.enqueue):
			model.Queue.Enqueue(item)
			return m.NewStatusMessage(statusMessageStyle(item.Title() + " added to queue"))
		}
	}
	return nil
}

func (sd SongDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	view.RenderSong(sd, w, m, index, listItem)
}
