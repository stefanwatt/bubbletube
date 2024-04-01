package model

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type YouTubePlaylistListResponse struct {
	Kind     string     `json:"kind"`
	Etag     string     `json:"etag"`
	Items    []Playlist `json:"items"`
	PageInfo PageInfo   `json:"pageInfo"`
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type Playlist struct {
	Kind    string  `json:"kind"`
	Etag    string  `json:"etag"`
	ID      string  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type Snippet struct {
	PublishedAt  time.Time  `json:"publishedAt"`
	Localized    Localized  `json:"localized"`
	ChannelID    string     `json:"channelId"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	ChannelTitle string     `json:"channelTitle"`
	Thumbnails   Thumbnails `json:"thumbnails"`
}

type Thumbnails struct {
	Default  Thumbnail `json:"default"`
	Medium   Thumbnail `json:"medium"`
	High     Thumbnail `json:"high"`
	Standard Thumbnail `json:"standard"`
	Maxres   Thumbnail `json:"maxres"`
}

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Localized struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func getPlaylists() YouTubePlaylistListResponse {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://youtube.googleapis.com/youtube/v3/playlists?part=snippet&maxResults=50&mine=true&key="+currentToken.AccessToken,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+currentToken.AccessToken)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var playlistResponse YouTubePlaylistListResponse
	err = json.Unmarshal(body, &playlistResponse)
	if err != nil {
		log.Fatal(err)
	}
	return playlistResponse
}

type SongsOfPlaylistResponse struct {
	Kind          string         `json:"kind"`
	Etag          string         `json:"etag"`
	NextPageToken string         `json:"nextPageToken"`
	Items         []PlaylistItem `json:"items"`
	PageInfo      PageInfo       `json:"pageInfo"`
}

type PlaylistItem struct {
	Kind    string      `json:"kind"`
	Etag    string      `json:"etag"`
	ID      string      `json:"id"`
	Snippet ItemSnippet `json:"snippet"`
}

type ItemSnippet struct {
	PublishedAt  time.Time  `json:"publishedAt"`
	ResourceID   ResourceID `json:"resourceId"`
	ChannelID    string     `json:"channelId"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	ChannelTitle string     `json:"channelTitle"`
	PlaylistID   string     `json:"playlistId"`
	Thumbnails   Thumbnails `json:"thumbnails"`
	Position     int        `json:"position"`
}

type ResourceID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

func getSongsOfPlaylist(playlistID string) SongsOfPlaylistResponse {
	client := &http.Client{}
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=50&playlistId=%s&key=%s",
		playlistID,
		currentToken.AccessToken,
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+currentToken.AccessToken)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res SongsOfPlaylistResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatal(err)
	}
	return res
}
