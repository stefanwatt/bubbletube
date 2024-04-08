package controller

import (
	model "bubbletube/model"
	view "bubbletube/view"

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
		key.WithHelp("ctrl+right", "Next Song"),
	),
	prevSong: key.NewBinding(
		key.WithKeys("ctrl+left"),
		key.WithHelp("ctrl+left", "Previous Song"),
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
	case tea.WindowSizeMsg:
		sc.Screen.WindowWidth = msg.Width
		sc.Screen.WindowHeight = msg.Height
		list := sc.Screen.CenterPanel.GetList()
		list.SetWidth(80)
		list.SetHeight(view.GetCenterPanelHeight(msg.Height) - 1)
		updatedList, cmd := list.Update(list)
		sc.Screen.CenterPanel.SetList(updatedList)
		return sc, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, globalKeys.noop):
			return sc, nil
		case key.Matches(msg, globalKeys.quit):
			sc.Screen.Quitting = true
			model.KillMpv()
			return sc, tea.Quit
		case key.Matches(msg, globalKeys.volumeDown):
			volume, err := model.VolumeDown()
			if err != nil {
				return sc, nil
			}
			sc.Screen.PlaybackControls.Volume = volume
			return sc, nil
		case key.Matches(msg, globalKeys.volumeUp):
			volume, err := model.VolumeUp()
			if err != nil {
				return sc, nil
			}
			sc.Screen.PlaybackControls.Volume = volume
			return sc, nil
		case key.Matches(msg, globalKeys.nextSong):
			sc.Screen.PlaybackControls = sc.Screen.QueuePanel.NextSong(sc.Screen.PlaybackControls)
			return sc, nil
		case key.Matches(msg, globalKeys.skipBackward):
			time_pos, time_remaining, err := model.SkipBackward()
			if err != nil {
				return sc, nil
			}
			sc.Screen.PlaybackControls.TimePos = time_pos
			sc.Screen.PlaybackControls.TimeRemaining = time_remaining
			return sc, nil
		case key.Matches(msg, globalKeys.skipForward):
			time_pos, time_remaining, err := model.SkipForward()
			if err != nil {
				return sc, nil
			}
			sc.Screen.PlaybackControls.TimePos = time_pos
			sc.Screen.PlaybackControls.TimeRemaining = time_remaining
			return sc, nil
		case key.Matches(msg, globalKeys.pause):
			sc.Screen.PlaybackControls.Playing = !model.TogglePlayback()
			return sc, nil
		}

	case model.MPVFloatValueChangedEvent:
		return MpvFloatValueUpdate(msg, sc)
	case model.MPVVoidValueChangedEvent:
		return MpvVoidValueUpdate(msg, sc)

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
