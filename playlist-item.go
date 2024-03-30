package main

type YTPlaylistItem struct {
	ID              string
	TitleText       string
	DescriptionText string
}

func NewYTPlaylistItem(id string, title string, description string) YTPlaylistItem {
	return YTPlaylistItem{
		ID:              id,
		TitleText:       title,
		DescriptionText: description,
	}
}

func (p YTPlaylistItem) FilterValue() string {
	return p.TitleText
}

func (p YTPlaylistItem) Title() string {
	return p.TitleText
}

func (p YTPlaylistItem) Description() string {
	return p.DescriptionText
}
