package dynamicprogramming

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"unicode"

	"github.com/golang/glog"
)

type Memo map[int]map[int]bool

func testEquality(memo Memo, source []rune, i int, target []rune, j int) bool {
	sourceRune := source[i]
	targetRune := target[j]
	// Once we elect to delete a lower-case letter from the source, we have to delete them all.
	match := sourceRune == targetRune || (unicode.ToUpper(sourceRune) == targetRune)
	return match && abbreviateFrom(memo, source, i+1, target, j+1)
}

func testDeletion(memo Memo, source []rune, i int, target []rune, j int) bool {
	sourceRune := source[i]
	// If we're skipping characters..
	return unicode.IsLower(sourceRune) && abbreviateFrom(memo, source, i+1, target, j)
}

func abbreviateFrom(memo Memo, source []rune, sourcePosition int, target []rune, targetPosition int) bool {
	// Have we used up the source string?
	if sourcePosition == len(source) {
		return targetPosition == len(target)
	}

	// Check memoization.
	if b, ok := memo[sourcePosition][targetPosition]; ok {
		return b
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
// See https://www.hackerrank.com/challenges/abbr/problem
func Abbreviate(source string, target string) bool {
	a := []rune(source)
	b := []rune(target)
	match := make(map[int]map[int]bool)
	for i := range a {
		match[i] = make(map[int]bool)
	}
	return abbreviateFrom(match, a, 0, b, 0)
}

func TestAbbreviation(t *testing.T) {
	tests := []struct {
		source string
		target string
		expect bool
	}{
		{"", "", true},
		{"abc", "", true},
		{"abc", "abc", true},
		{"ABC", "ABC", true},
		{"abCde", "C", true},
		// I'm puzzled about this one.
		{"abCd", "Cd", false},
		{"abcdE", "E", true},
	}
	for _, test := range tests {
		if test.expect {
			if !Abbreviate(test.source, test.target) {
				t.Errorf("%s should match %s", test.source, test.target)
			}
		} else {
			if Abbreviate(test.source, test.target) {
				t.Errorf("%s should not match %s", test.source, test.target)
			}
		}
	}
}

func TestAbbreviationTestCases(t *testing.T) {
	name := "dynamicprogramming/abbreviation-input12.txt"
	path, _ := filepath.Abs(name)
	file, e := os.Open(path)
	if e != nil {
		glog.Error(e)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	te := scanner.Text()
	n, _ := strconv.Atoi(te)
	for i := 0; i < n; i++ {
		scanner.Scan()
		source := scanner.Text()
		scanner.Scan()
		target := scanner.Text()
		result := Abbreviate(source, target)
		fmt.Printf("Result of %s... / %s...: %t\n", source[0:10], target[0:10], result)
	}
}

