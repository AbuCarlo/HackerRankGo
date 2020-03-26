package dictionaries

// TwoStrings checks for the simplest common substring,
// i.e. one common character.
func TwoStrings(s string, t string) bool {
	d := map[rune]bool{}
	for _, c := range s {
		d[c] = true
	}
	for _, c := range t {
		if d[c] {
			return true
		}
	}
	return false
}
