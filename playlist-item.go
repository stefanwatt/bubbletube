package main

type YTPlaylist struct {
	ID              string
	TitleText       string
	DescriptionText string
}

func NewYTPlaylistItem(id string, title string, description string) YTPlaylist {
	return YTPlaylist{
		ID:              id,
		TitleText:       title,
		DescriptionText: description,
	}
}

func (p YTPlaylist) FilterValue() string {
	return p.TitleText
}

func (p YTPlaylist) Title() string {
	return p.TitleText
}

func (p YTPlaylist) Description() string {
	return p.DescriptionText
}
