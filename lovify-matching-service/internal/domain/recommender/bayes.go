package recommender

import "math"

// cosSim calcula la similitud coseno entre dos vectores esparcidos
func cosSim(a, b map[string]float64) float64 {
	var dot, na, nb float64
	for g, va := range a {
		if vb, ok := b[g]; ok {
			dot += va * vb
		}
		na += va * va
	}
	for _, vb := range b {
		nb += vb * vb
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}
