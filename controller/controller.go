package controller

import (
	model "bubbletube/model"
	view "bubbletube/view"
	"math"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
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

type globalKeymap struct {
	quit         key.Binding
	volumeDown   key.Binding
	volumeUp     key.Binding
	nextSong     key.Binding
	prevSong     key.Binding
	skipForward  key.Binding
	skipBackward key.Binding
	pause        key.Binding
	noop         key.Binding
}

var globalKeys = &globalKeymap{
	noop: key.NewBinding(
		key.WithKeys("q"),
	),
	quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "Exit BubbleTube"),
	),
	volumeDown: key.NewBinding(
		key.WithKeys("ctrl+down"),
		key.WithHelp("ctrl+down", "volume down"),
	),
	volumeUp: key.NewBinding(
		key.WithKeys("ctrl+up"),
		key.WithHelp("ctrl+up", "volume up"),
	),
	nextSong: key.NewBinding(
		key.WithKeys("ctrl+right"),
		key.WithHelp("ctrl+up", "Next Song"),
	),
	prevSong: key.NewBinding(
		key.WithKeys("ctrl+left"),
		key.WithHelp("ctrl+up", "Previous Song"),
	),
	skipForward: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("right", "Skip Forward"),
	),
	skipBackward: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("left", "Skip Backward"),
	),
	pause: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Pause/Play"),
	),
}

func (sc *ScreenController) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, globalKeys.quit):
			sc.Screen.Quitting = true
			model.KillMpv()
			return sc, tea.Quit
		case key.Matches(msg, globalKeys.volumeDown):
			model.VolumeDown()
			sc.Screen.PlaybackControls.Volume = sc.Screen.PlaybackControls.Volume - 5
			return sc, nil
		case key.Matches(msg, globalKeys.volumeUp):
			model.VolumeUp()
			sc.Screen.PlaybackControls.Volume = sc.Screen.PlaybackControls.Volume + 5
			return sc, nil
		case key.Matches(msg, globalKeys.nextSong):
			sc.Screen.PlaybackControls = model.NextSong(sc.Screen.PlaybackControls)
			return sc, nil
		case key.Matches(msg, globalKeys.skipBackward):
			model.SkipBackward()
			sc.Screen.PlaybackControls.TimePos = sc.Screen.PlaybackControls.TimePos - 10
			return sc, nil
		case key.Matches(msg, globalKeys.skipForward):
			model.SkipForward()
			sc.Screen.PlaybackControls.TimePos = sc.Screen.PlaybackControls.TimePos + 10
			return sc, nil
		case key.Matches(msg, globalKeys.pause):
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
