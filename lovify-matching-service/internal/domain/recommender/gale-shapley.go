package recommender

func isFree(matched map[string]string, id string) bool {
	_, taken := matched[id]
	return !taken
}
func StableMatch(preferences map[string][]ScoredUser) map[string]string {
	matched := map[string]string{}
	next := map[string]int{}

	// iteramos sobre todos los proponentes iniciales
	queue := []string{}
	for id := range preferences {
		queue = append(queue, id)
	}

	for len(queue) > 0 {
		proposer := queue[0]
		queue = queue[1:]

		if next[proposer] >= len(preferences[proposer]) {
			continue
		}
		target := preferences[proposer][next[proposer]].OtherID
		next[proposer]++

		if isFree(matched, target) {
			matched[proposer] = target
			matched[target] = proposer
			continue
		}
		// target ya tiene pareja: ¿prefiere al nuevo?
		current := matched[target]
		if prefers(target, proposer, current, preferences) {
			// rompe con current
			delete(matched, current)
			queue = append(queue, current)
			matched[proposer] = target
			matched[target] = proposer
		} else {
			queue = append(queue, proposer)
		}
	}
	// devuelve solo una entrada por pareja
	final := map[string]string{}
	for a, b := range matched {
		if _, ok := final[b]; !ok {
			final[a] = b
		}
	}
	return final
}

// prefers devuelve true si t prefiere a new en lugar de current
func prefers(target, newID, currentID string, preferences map[string][]ScoredUser) bool {
	rank := map[string]int{}
	for i, s := range preferences[target] {
		rank[s.OtherID] = i
	}
	return rank[newID] < rank[currentID]
}
