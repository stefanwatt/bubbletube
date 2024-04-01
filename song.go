package main

type SongItem struct {
	ID        string
	VideoID   string
	TitleText string
	Artist    string
}

func NewSongItem(id string, title string, videoID string) SongItem {
	return SongItem{
		ID:        id,
		VideoID:   videoID,
		TitleText: title,
	}
}

func (p SongItem) FilterValue() string {
	return p.TitleText
}

func (p SongItem) Title() string {
	return p.TitleText
}
