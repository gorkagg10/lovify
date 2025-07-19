package recommender

import "time"

type User struct {
	ID                       string
	Email                    string
	Birthday                 time.Time
	Name                     string
	Gender                   string
	SexualOrientation        string
	Description              string
	ConnectedToMusicProvider bool
	MusicProviderInfo        *MusicProviderData
}

type MusicProviderData struct {
	TopTracks  []Track
	TopArtists []Artist
}

type Album struct {
	Name      string
	AlbumType string
	Image     *Image
}

type Artist struct {
	Name   string
	Genres []string
	Image  *Image
}

type Image struct {
	Url    string
	Height int
	Width  int
}

type Track struct {
	Name    string
	Album   *Album
	Artists []string
}
