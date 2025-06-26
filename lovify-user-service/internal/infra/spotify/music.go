package spotify

type TopTracks struct {
	Tracks []Tracks `json:"items"`
}

type Tracks struct {
	Name    string   `json:"name"`
	Album   Album    `json:"album"`
	Artists []Artist `json:"artists"`
}

type Album struct {
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Images []Image `json:"images"`
}

type Artist struct {
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
	Images []Image  `json:"images"`
}

type TopArtists struct {
	Artists []Artist `json:"items"`
}

type Image struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}
