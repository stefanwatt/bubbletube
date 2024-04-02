package model

var Queue PlaybackQueue

type PlaybackQueue []SongItem

func (q *PlaybackQueue) Enqueue(item SongItem) {
	*q = append(*q, item)
}

func (q *PlaybackQueue) Dequeue() (*SongItem, bool) {
	if len(*q) == 0 {
		return nil, false
	}
	item := (*q)[0]
	*q = (*q)[1:]
	return &item, true
}

func (q *PlaybackQueue) GetQueue() []SongItem {
	return *q
}

func NextSong(playbackControls PlaybackControls) PlaybackControls {
	nextSong, ok := Queue.Dequeue()
	if ok {
		playbackControls.PlayingSong = nextSong
		SelectSong(*nextSong)
	}
	return playbackControls
}
