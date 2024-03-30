package main

type SongItem struct {
	ID        string
	TitleText string
}

func NewSongItem(id string, title string, description string) SongItem {
	return SongItem{
		ID:        id,
		TitleText: title,
	}
}

func (p SongItem) FilterValue() string {
	return p.TitleText
}

func (p SongItem) Title() string {
	return p.TitleText
}
