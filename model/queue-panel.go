package model

import (
	"bubbletube/utils"

	"github.com/charmbracelet/bubbles/list"
)

type QueuePanel struct {
	Waitlist list.Model
	Playlist list.Model
	Choice   SongItem
}

func (q *QueuePanel) Enqueue(item SongItem) {
	q.Waitlist.InsertItem(0, item)
}

func (q *QueuePanel) Dequeue() (*SongItem, bool) {
	items := q.Waitlist.Items()
	if len(items) == 0 {
		if len(q.Playlist.Items()) == 0 {
			return nil, false
		}
		item := q.Playlist.Items()[0]
		songItem, ok := item.(SongItem)
		if !ok {
			return nil, false
		}
		q.Playlist.RemoveItem(0)
		return &songItem, true
	}
	item := q.Waitlist.Items()[0]
	songItem, ok := item.(SongItem)
	if !ok {
		return nil, false
	}
	q.Waitlist.RemoveItem(0)
	return &songItem, true
}

func (q *QueuePanel) GetQueue() []SongItem {
	items := q.Waitlist.Items()
	songItems := utils.MapArray(items, func(item list.Item) SongItem {
		songItem, _ := item.(SongItem)
		return songItem
	})
	return songItems
}

func (q *QueuePanel) NextSong(playbackControls PlaybackControls) PlaybackControls {
	nextSong, ok := q.Dequeue()
	if ok {
		playbackControls.PlayingSong = nextSong
		SelectSong(*nextSong)
	}
	return playbackControls
}
