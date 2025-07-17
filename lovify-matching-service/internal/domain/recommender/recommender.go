package recommender

import "sort"

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

func compatible(userA, userB User) bool {
	if userA.SexualOrientation == "HETEROSEXUAL" && userB.SexualOrientation == "HETEROSEXUAL" {
		return userA.Gender != userB.Gender
	}
	if userA.SexualOrientation == "HOMOSEXUAL" && userB.SexualOrientation == "HOMOSEXUAL" {
		return userA.Gender == userB.Gender
	}
	return true
}

// BuildPreferences creates for each user an ordered list of compatible candidates.
func BuildPreferences(users []User, minScore float64) map[string][]ScoredUser {
	preferences := make(map[string][]ScoredUser, len(users))

	// vector precalculation
	vectors := map[string]map[string]float64{}
	for _, user := range users {
		if user.MusicProviderInfo != nil {
			vectors[user.Email] = BuildGenreVector(user.MusicProviderInfo.TopArtists)
		} else {
			vectors[user.Email] = map[string]float64{}
		}
	}

	for _, user := range users {
		for _, vectorUser := range users {
			if user.Email == vectorUser.Email || !compatible(user, vectorUser) {
				continue
			}
			score := CosSim(vectors[user.Email], vectors[vectorUser.Email])
			if score >= minScore {
				preferences[user.Email] = append(preferences[user.Email], ScoredUser{
					OtherID: vectorUser.Email,
					Score:   score,
				})
			}
		}
		sort.Slice(preferences[user.Email], func(i, j int) bool {
			return preferences[user.Email][i].Score > preferences[user.Email][j].Score
		})
	}
	return preferences
}
