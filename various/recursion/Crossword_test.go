package recursion

import (
	"math"
	"strings"
	"testing"
)

const Boundary = 10
const Size = 10
const FillCharacter = '-'

type Word struct {
	word        string
	across      bool
	row, column int
}

type Slot struct {
	across      bool
	row, column, length int
}

type Xword []Word

func populateGrid() [][]rune {
	template := make([]rune, Boundary)
	for i, _ := range template {
		template[i] = FillCharacter
	}
	rows := [][]rune{}
	for i := 0; i < Boundary; i++ {
		row := make([]rune, len(template))
		copy(row, template)
		rows = append(rows, row)
	}
	return rows
}

func (xword Xword) render() string {
	grid := populateGrid()

	upperLeft, _ := findBoundaries(xword)

	for _, word := range xword {
		for j, c := range word.word {
			if word.across {
				grid[word.row-upperLeft.row][word.column-upperLeft.column+j] = c
			} else {
				grid[word.row-upperLeft.row+j][word.column-upperLeft.column] = c
			}
		}
	}

	result := make([]string, 10)
	for i := 0; i < len(grid); i++ {
		result[i] = string(grid[i])
	}
	return strings.Join(result, "\n")
}

func overlap(w, x Word) bool {
	if w.across != x.across {
		return false
	}

	if w.across {
		if w.row != x.row {
			return false
		}

		if w.column <= x.column {
			return w.column+len(w.word)-1 >= x.column
		} else {
			return x.column+len(x.word)-1 >= w.column
		}
	}

	if w.column != x.column {
		return false
	}

	if w.row <= x.row {
		return w.row+len(w.word) >= x.row
	} else {
		return x.row+len(x.word) >= w.row
	}
}

type Coordinate struct {
	row, column int
}

func findBoundaries(words []Word) (Coordinate, Coordinate) {
	firstRow := math.MaxInt
	lastRow := math.MinInt
	firstColumn := math.MaxInt
	lastColumn := math.MinInt

	for _, w := range words {
		firstRow = min(firstRow, w.row)
		if w.across {
			lastRow = max(lastRow, w.row+len(w.word)-1)
		} else {
			lastRow = max(lastRow, w.row)
		}
		firstColumn = min(firstColumn, w.column)
		if w.across {
			lastColumn = max(lastColumn, w.column)
		} else {
			lastColumn = max(lastColumn, w.column+len(w.word)-1)
		}
	}

	return Coordinate{firstRow, firstColumn}, Coordinate{lastRow, lastColumn}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func isAlongsideHorizontal(w, x Word) bool {
	if abs(w.row-x.row) > 1 {
		return false
	}
	if w.column > x.column {
		return isAlongsideHorizontal(x, w)
	}

	wLast := w.column + len(w.word) - 1
	xLast := x.column + len(x.word) - 1
	return w.column >= x.column && w.column <= xLast || wLast >= x.column && wLast <= xLast
}

func isAlongsideVertical(w, x Word) bool {
	if abs(w.column-x.column) > 1 {
		return false
	}
	if w.row > x.row {
		return isAlongsideVertical(x, w)
	}

	wLast := w.row + len(w.word) - 1
	xLast := x.row + len(x.word) - 1
	return w.row >= x.row && w.row <= xLast || wLast >= x.row && wLast <= xLast
}

func isAlongside(w, x Word) bool {
	if w.across != x.across {
		return false
	}

	if w.across {
		return isAlongsideHorizontal(w, x)
	}
	return isAlongsideVertical(w, x)
}

func isWithinBoundaries(xword []Word, w Word) bool {
	// Too much copying here.
	newXword := append(xword, w)

	upperLeft, lowerRight := findBoundaries(newXword)
	return lowerRight.column-upperLeft.column+1 > Boundary || lowerRight.row-upperLeft.row+1 > Boundary
}

func findCrossings(w Word, s string) []Word {
	var crossings []Word
	if w.across {
		// Pretend that w starts at (0, 0). The math is easier.
		for c := 0; c < len(w.word); c++ {
			for r := -len(s) + 1; r < 1; r++ {
				if w.word[c] == s[r+len(s)-1] {
					crossings = append(crossings, Word{s, false, r + w.row, c + w.column})
				}
			}
		}
	} else {
		// Pretend that w starts at (0, 0)
		for r := 0; r < len(w.word); r++ {
			for c := -len(s) + 1; c < 1; c++ {
				if w.word[r] == s[c+len(s)-1] {
					crossings = append(crossings, Word{s, true, c + w.column, r + w.row})
				}
			}
		}
	}

	return crossings
}

func allowed(xword []Word, word Word) bool {
	if !isWithinBoundaries(xword, word) {
		return false
	}

	for _, x := range xword {
		if isAlongside(word, x) {
			return false
		}
		if overlap(x, word) {
			return false
		}
	}

	return true
}

func recurse(xword []Word, words []string) ([]Word, bool) {
	if len(words) == 0 {
		return xword, true
	}

	for _, w := range words {
		xings := findCrossings(xword[len(xword)-1], w)
		for _, xing := range xings {
			if !allowed(xword, xing) {
				continue
			}
			next := append(xword, xing)
			if blah, ok := recurse(next, words[1:]); ok {
				return blah, true
			}
		}
	}

	return nil, false
}

func createCrossword(words []string) []Word {
	if len(words) == 0 {
		return []Word{}
	}
	across := Word{words[0], true, 0, 0}
	if xword, ok := recurse([]Word{across}, words[1:]); ok {
		return xword
	}
	down := Word{words[0], false, 0, 0}
	if xword, ok := recurse([]Word{down}, words[1:]); ok {
		return xword
	}
	return nil
}

func TestString(t *testing.T) {
	xword := Xword{}
	s := xword.render()
	t.Logf("%s", s)

	alien := Xword{{"ALIEN", true, 0, 0}, {"ALIEN", false, 0, 0}}
	s = alien.render()
	t.Logf("%s", s)

	animals := Xword{{"ALIEN", true, 0, 0}, {"ANIMAL", false, -1, 4}}
	s = animals.render()
	t.Logf("%s", s)
}

func TestCrossings(t *testing.T) {
	// For any word with no repeated letters,
	// the number of crossings must equal the
	// number of letters.
	word := Word{"ALIEN", true, 0, 0}
	crossings := findCrossings(word, "ALIEN")
	t.Logf("Found %v", crossings)

	badWord := Word{"SYRINX", true, 0, 0}
	noCrossings := findCrossings(badWord, "MOPE")
	if len(noCrossings) > 0 {
		t.Error("Should be 0")
	}
}

func TestCollisions(t *testing.T) {
	type TestCase struct {
		w        Word
		x        Word
		expected bool
	}
	table := []TestCase{
		{Word{"GWALIOR", true, 0, 0}, Word{"GWALIOR", true, 0, 0}, true},
		{Word{"GWALIOR", true, 0, 0}, Word{"GWALIOR", false, 0, 0}, false},
		{Word{"GWALIOR", true, 0, 0}, Word{"GWALIOR", true, 0, 1}, true},
		{Word{"GWALIOR", true, 0, 0}, Word{"GWALIOR", true, 0, 6}, true},
		{Word{"GWALIOR", true, 0, 0}, Word{"GWALIOR", true, 0, 7}, false},
	}
	for i, test := range table {
		if overlap(test.w, test.x) != test.expected {
			t.Errorf("Row %d: %v and %v collision expected %v", i, test.w, test.x, test.expected)
		}

		// Test the commutation!
		if overlap(test.x, test.w) != test.expected {
			t.Errorf("Row %d: %v and %v collision (inverted) expected %v", i, test.w, test.x, test.expected)
		}
	}
}

func findSlots(xword []string) []Slot {
	slots := []Slot{}
	for row, s := range xword {
		column := 0
		for column < len(s) {
			for ; column < Size && s[column] == '+'; column++ {
				//
			}
			if column == Size {
				continue
			}
			lastColumn := column
			for ; lastColumn < Size && s[lastColumn] == '-'; lastColumn++ {
				//
			}
			if lastColumn > column + 1 {
				slot := Slot{true, row, column, lastColumn - column}
				slots = append(slots, slot)
				column = lastColumn
			} else {
				column++
			}
		}
	}

	for column := range Size {
		row := 0
		for row < Size {
			for ; row < Size && xword[row][column] == '+'; row++ {
				//
			}
			if row ==  Size {
				continue
			}
			lastRow := row
			for ; lastRow < Size && xword[lastRow][column] == '-'; lastRow++ {
				//
			}
			if lastRow > row + 1 {
				slot := Slot{false, row, column, lastRow - row}
				slots = append(slots, slot)
				row = lastRow
			} else {
				row++
			}
		}
	}

	return slots
}

func TestSamples(t *testing.T) {
	table := [][]string{
		{
			"++++++++++",
			"+------+++",
			"+++-++++++",
			"+++-++++++",
			"+++-----++",
			"+++-++-+++",
			"++++++-+++",
			"++++++-+++",
			"++++++-+++",
			"++++++++++"},
	}
	for _, row := range table {
		answer := findSlots(row)
		t.Logf("%v", answer)
	}
}
