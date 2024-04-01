package model

type Playlist struct {
	ID              string
	TitleText       string
	DescriptionText string
}

func NewYTPlaylistItem(id string, title string, description string) Playlist {
	return Playlist{
		ID:              id,
		TitleText:       title,
		DescriptionText: description,
	}
}

func (p Playlist) FilterValue() string {
	return p.TitleText
}

func (p Playlist) Title() string {
	return p.TitleText
}

func (p Playlist) Description() string {
	return p.DescriptionText
}
