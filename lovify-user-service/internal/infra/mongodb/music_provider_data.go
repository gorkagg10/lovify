package mongodb

type MusicProviderData struct {
	UserID    string   `bson:"user_id"`
	TopTracks []Track  `bson:"top_tracks"`
	TopArtist []Artist `bson:"top_artist"`
}

func NewMusicProviderData(
	userID string,
	topTracks []Track,
	topArtist []Artist,
) *MusicProviderData {
	return &MusicProviderData{
		UserID:    userID,
		TopTracks: topTracks,
		TopArtist: topArtist,
	}
}

type Track struct {
	Name    string   `bson:"name"`
	Album   *Album   `bson:"album"`
	Artists []string `bson:"artists"`
}

func NewTrack(name string, album *Album, artists []string) *Track {
	return &Track{
		Name:    name,
		Album:   album,
		Artists: artists,
	}
}

type Album struct {
	Name  string `bson:"name"`
	Type  string `bson:"type"`
	Cover *Image `bson:"cover"`
}

func NewAlbum(name string, albumType string, cover *Image) *Album {
	return &Album{
		Name:  name,
		Type:  albumType,
		Cover: cover,
	}
}

type Artist struct {
	Name   string   `bson:"name"`
	Genres []string `bson:"genres"`
	Image  *Image   `bson:"image"`
}

func NewArtist(name string, genres []string, image *Image) *Artist {
	return &Artist{
		Name:   name,
		Genres: genres,
		Image:  image,
	}
}

type Image struct {
	Url    string `bson:"url"`
	Height int    `bson:"height"`
	Width  int    `bson:"width"`
}

func NewImage(url string, height int, width int) *Image {
	return &Image{
		Url:    url,
		Height: height,
		Width:  width,
	}
}
