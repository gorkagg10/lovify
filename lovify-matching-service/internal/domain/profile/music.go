package profile

type MusicProviderData struct {
	topTracks  []Track
	topArtists []Artist
}

func NewMusicProviderData(topTracks []Track, topArtists []Artist) *MusicProviderData {
	return &MusicProviderData{
		topTracks:  topTracks,
		topArtists: topArtists,
	}
}

func (m *MusicProviderData) TopTracks() []Track {
	return m.topTracks
}

func (m *MusicProviderData) TopArtists() []Artist {
	return m.topArtists
}

type Track struct {
	name    string
	album   *Album
	artists []string
}

func NewTrack(name string, album *Album, artists []string) *Track {
	return &Track{
		name:    name,
		album:   album,
		artists: artists,
	}
}

func (t *Track) Name() string {
	return t.name
}

func (t *Track) Album() *Album {
	return t.album
}

func (t *Track) Artists() []string {
	return t.artists
}

type Album struct {
	name      string
	albumType string
	image     *Image
}

func NewAlbum(name string, albumType string, image *Image) *Album {
	return &Album{
		name:      name,
		albumType: albumType,
		image:     image,
	}
}

func (a *Album) Name() string {
	return a.name
}

func (a *Album) Type() string {
	return a.albumType
}

func (a *Album) Image() *Image {
	return a.image
}

type Artist struct {
	name   string
	genres []string
	image  *Image
}

func NewArtist(name string, genres []string, image *Image) *Artist {
	return &Artist{
		name:   name,
		genres: genres,
		image:  image,
	}
}

func (a *Artist) Name() string {
	return a.name
}

func (a *Artist) Genres() []string {
	return a.genres
}

func (a *Artist) Image() *Image {
	return a.image
}

type Image struct {
	url    string
	height int
	width  int
}

func NewImage(url string, height int, width int) *Image {
	return &Image{
		url:    url,
		height: height,
		width:  width,
	}
}

func (i *Image) URL() string {
	return i.url
}

func (i *Image) Height() int {
	return i.height
}

func (i *Image) Width() int {
	return i.width
}
