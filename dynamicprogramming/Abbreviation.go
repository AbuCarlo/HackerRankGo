package dynamicprogramming

import (
	"unicode"
)

func testEquality(memo [][]bool, source []rune, i int, target []rune, j int) bool {
	sourceRune := source[i]
	targetRune := target[j]
	// Once we elect to delete a lower-case letter from the source, we have to delete them all.
	match := sourceRune == targetRune || (unicode.ToUpper(sourceRune) == targetRune)
	return match && abbreviateFrom(memo, source, i+1, target, j+1)
}

func testDeletion(memo [][]bool, source []rune, i int, target []rune, j int) bool {
	sourceRune := source[i]
	// If we're skipping characters..
	return unicode.IsLower(sourceRune) && abbreviateFrom(memo, source, i+1, target, j)
}

func abbreviateFrom(memo [][]bool, source []rune, sourcePosition int, target []rune, targetPosition int) bool {
	// Have we used up the source string?
	if sourcePosition == len(source) {
		return targetPosition == len(target)
	}

	// Check memoization.
	if len(memo[sourcePosition]) > targetPosition {
		return memo[sourcePosition][targetPosition]
	}

	if targetPosition >= len(memo[sourcePosition]) {
		appendage := make([]bool, targetPosition-len(memo[sourcePosition])+1)
		memo[sourcePosition] = append(memo[sourcePosition], appendage...)
	}

	// Don't short-circuit!
	matchWithDeletion := testDeletion(memo, source, sourcePosition, target, targetPosition)
	matchWithEquality := targetPosition < len(target) && testEquality(memo, source, sourcePosition, target, targetPosition)
	result := matchWithDeletion || matchWithEquality
	// fmt.Printf("Checking %s (%d) against %s (%d) with %t\n", string(source[sourcePosition:]), sourcePosition, string(target[targetPosition:]), targetPosition, result)
	memo[sourcePosition][targetPosition] = result
	return result
}

// Abbreviate is a version of "longest common subsequence".
// See https://www.hackerrank.com/challenges/abbr/problem?h_l=interview&playlist_slugs%5B%5D=interview-preparation-kit&playlist_slugs%5B%5D=dynamic-programming.
func Abbreviate(source string, target string) bool {
	a := []rune(source)
	b := []rune(target)
	match := make([][]bool, len(a))
	return abbreviateFrom(match, a, 0, b, 0)
}
