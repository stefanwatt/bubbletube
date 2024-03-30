package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	items             []list.Item
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	ytItem, ok := listItem.(YTPlaylistItem)
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

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "p":
			TogglePlayback()
			return m, nil

		case "enter":
			i, ok := m.list.SelectedItem().(YTPlaylistItem)
			if ok {
				m.choice = string(i.ID)
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	m.choice = ""
	return m, cmd
}

func (m model) View() string {
	item, err := Find(items, func(i list.Item) bool {
		return i.(YTPlaylistItem).ID == m.choice
	})
	if m.choice != "" && err == nil {
		fmt.Println("choice: " + m.choice)
		playlistItem := item.(YTPlaylistItem)
		songs := getPlaylistItems(playlistItem.ID)

		PlaySong(songs.Items[0])
		return quitTextStyle.Render(fmt.Sprintf("Playing: %s", songs.Items[0].Snippet.Title))
	}
	if m.quitting {
		return quitTextStyle.Render("Bye")
	}
	return "\n" + m.list.View()
}

func Find[T any](arr []T, f func(T) bool) (T, error) {
	var zero T
	for _, value := range arr {
		if f(value) {
			return value, nil
		}
	}
	return zero, fmt.Errorf("no match found")
}

func MapArray[T any, U any](arr []T, f func(T) U) []U {
	var result []U
	for _, value := range arr {
		result = append(result, f(value))
	}
	return result
}

func main() {
	err := godotenv.Load() // Load .env file from the current directory
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	items = []list.Item{}
	res := getPlaylists()
	for _, playlist := range res.Items {
		playlistItem := YTPlaylistItem{
			ID:              playlist.ID,
			TitleText:       playlist.Snippet.Title,
			DescriptionText: playlist.Snippet.Description,
		}
		items = append(items, playlistItem)
	}

	const defaultWidth = 100

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = fmt.Sprintf("%d playlists found", len(items))
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
