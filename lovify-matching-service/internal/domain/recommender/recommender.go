package recommender

import (
	"sort"
	"strings"
)

type ScoredUser struct {
	OtherID string
	Score   float64
}

type GenreKeyValue struct {
	Genre string
	Freq  float64
}

// BuildGenreVector devuelve genero -> frecuencia relativa dentro de los artistas favoritos.
// Normaliza para que ∑score = 1.
func BuildGenreVector(artists []Artist) map[string]float64 {
	genreCounter := map[string]float64{}
	for _, artist := range artists {
		for _, genre := range artist.Genres {
			genreCounter[genre]++
		}
	}

	// normalize
	var sum float64
	for _, v := range genreCounter {
		sum += v
	}
	if sum == 0 {
		return genreCounter
	}
	for genre := range genreCounter {
		genreCounter[genre] /= sum
	}
	return genreCounter
}

// Devuelve un mapa genero → frecuencia normalizada
func GenreVector(u *User, topN int) map[string]float64 {
	if u.MusicProviderInfo == nil {
		return nil
	}
	count := map[string]float64{}
	for _, art := range u.MusicProviderInfo.TopArtists {
		for _, g := range art.Genres {
			count[strings.ToLower(g)]++
		}
	}
	// Selecciona los N géneros principales
	type kv struct {
		g string
		f float64
	}
	var kvs []kv
	for g, f := range count {
		kvs = append(kvs, kv{g, f})
	}
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].f > kvs[j].f })
	if len(kvs) > topN {
		kvs = kvs[:topN]
	}
	vec := map[string]float64{}
	var sum float64
	for _, kv := range kvs {
		sum += kv.f
	}
	for _, kv := range kvs {
		vec[kv.g] = kv.f / sum
	}
	return vec
}
