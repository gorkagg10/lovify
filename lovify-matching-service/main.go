package main

import (
	"fmt"
	"github.com/gorkagg10/lovify/lovify-matching-service/internal/domain/recommender"
)

func main() {
	users := []recommender.User{
		{
			Email:                    "alice@mail.com",
			Name:                     "Alice",
			Gender:                   "F",
			SexualOrientation:        "HETEROSEXUAL",
			ConnectedToMusicProvider: true,
			MusicProviderInfo: &recommender.MusicProviderData{
				TopArtists: []recommender.Artist{
					{Name: "ArtistA", Genres: []string{"pop", "rock"}},
					{Name: "ArtistB", Genres: []string{"pop"}},
				},
			},
		},
		{
			Email:                    "bob@mail.com",
			Name:                     "Bob",
			Gender:                   "M",
			SexualOrientation:        "HETEROSEXUAL",
			ConnectedToMusicProvider: true,
			MusicProviderInfo: &recommender.MusicProviderData{
				TopArtists: []recommender.Artist{
					{Name: "ArtistC", Genres: []string{"pop", "dance"}},
					{Name: "ArtistD", Genres: []string{"pop"}},
				},
			},
		},
		{
			Email:                    "carol@mail.com",
			Name:                     "Carol",
			Gender:                   "F",
			SexualOrientation:        "HETEROSEXUAL",
			ConnectedToMusicProvider: true,
			MusicProviderInfo: &recommender.MusicProviderData{
				TopArtists: []recommender.Artist{
					{Name: "ArtistE", Genres: []string{"rock", "metal"}},
					{Name: "ArtistF", Genres: []string{"metal"}},
				},
			},
		},
		{
			Email:                    "dave@mail.com",
			Name:                     "Dave",
			Gender:                   "M",
			SexualOrientation:        "HETEROSEXUAL",
			ConnectedToMusicProvider: true,
			MusicProviderInfo: &recommender.MusicProviderData{
				TopArtists: []recommender.Artist{
					{Name: "ArtistG", Genres: []string{"rock"}},
					{Name: "ArtistH", Genres: []string{"rock", "metal"}},
				},
			},
		},
	}

	minScore := 0.1 // umbral de compatibilidad
	preferences := recommender.BuildPreferences(users, minScore)
	possibleMatches := recommender.StableMatch(preferences)

	fmt.Println(possibleMatches)
	for u, v := range possibleMatches {
		fmt.Printf("  %s ↔ %s\n", u, v)
	}
}
