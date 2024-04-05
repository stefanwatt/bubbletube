package controller

import (
	model "bubbletube/model"
	view "bubbletube/view"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
			// model.Queue.Enqueue(item)
			// TODO: can we not handle this here?
			return m.NewStatusMessage(statusMessageStyle(item.Title() + " added to queue"))
		}
	}
	return nil
}

func (sd SongDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	view.RenderSong(sd, w, m, index, listItem)
}
