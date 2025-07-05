package spotify

import (
	"encoding/json"
	"fmt"
	"github.com/gorkagg10/lovify/lovify-user-service/internal/domain/profile"
	"io"
	"net/http"
)

const (
	TopTracksEndpoint  = "https://api.spotify.com/v1/me/top/tracks"
	TopArtistsEndpoint = "https://api.spotify.com/v1/me/top/artists"
)

type MusicProviderRepository struct {
	SpotifyClient *http.Client
}

func NewMusicProviderRepository(spotifyClient *http.Client) *MusicProviderRepository {
	return &MusicProviderRepository{
		SpotifyClient: spotifyClient,
	}
}

func (m *MusicProviderRepository) GetTopTracks(token string) ([]profile.Track, error) {
	req, err := http.NewRequest("GET", TopTracksEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("time_range", "short_term")

	resp, err := m.SpotifyClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Spotify API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tracks TopTracks
	if err = json.Unmarshal(body, &tracks); err != nil {
		return nil, err
	}

	topTracks := make([]profile.Track, len(tracks.Tracks))
	for i, track := range tracks.Tracks {
		artists := make([]string, len(tracks.Tracks[i].Artists))
		for j, artist := range tracks.Tracks[i].Artists {
			artists[j] = artist.Name
		}
		albumCover := profile.NewImage(track.Album.Images[0].Url, track.Album.Images[0].Height, track.Album.Images[0].Width)
		album := profile.NewAlbum(track.Album.Name, track.Album.Type, albumCover)
		topTracks[i] = *profile.NewTrack(
			track.Name,
			album,
			artists,
		)
	}
	return topTracks, nil
}

func (m *MusicProviderRepository) GetTopArtists(token string) ([]profile.Artist, error) {
	req, err := http.NewRequest("GET", TopArtistsEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("time_range", "short_term")

	resp, err := m.SpotifyClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Spotify API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var topArtists TopArtists
	if err = json.Unmarshal(body, &topArtists); err != nil {
		return nil, err
	}

	artists := make([]profile.Artist, len(topArtists.Artists))
	for i, _ := range artists {
		artists[i] = *profile.NewArtist(
			topArtists.Artists[i].Name,
			topArtists.Artists[i].Genres,
			profile.NewImage(
				topArtists.Artists[i].Images[0].Url,
				topArtists.Artists[i].Images[0].Height,
				topArtists.Artists[i].Images[0].Width,
			),
		)
	}
	return artists, nil
}
