package recursion

import (
	"math/rand"
	"strings"
	"testing"
)

const Size = 10

type Slot struct {
	across              bool
	row, column, length int
	word                string
}

func populateGrid() [][]rune {
	template := make([]rune, Size)
	for i, _ := range template {
		template[i] = '+'
	}
	rows := [][]rune{}
	for i := 0; i < Size; i++ {
		row := make([]rune, len(template))
		copy(row, template)
		rows = append(rows, row)
	}
	return rows
}

func render(answer []Slot) []string {
	grid := populateGrid()

	for _, o := range answer {
		for j, c := range o.word {
			if o.across {
				grid[o.row][o.column+j] = c
			} else {
				grid[o.row+j][o.column] = c
			}
		}
	}

	result := make([]string, 10)
	for i := 0; i < len(grid); i++ {
		result[i] = string(grid[i])
	}
	
	return result
}

func findSlots(xword []string) []Slot {
	slots := []Slot{}
	for row, s := range xword {
		column := 0
		for column < Size {
			for ; column < Size && s[column] == '+'; column++ {
				//
			}
			if column == Size {
				break
			}
			lastColumn := column
			for ; lastColumn < Size && s[lastColumn] == '-'; lastColumn++ {
				//
			}
			if lastColumn > column+1 {
				slot := Slot{true, row, column, lastColumn - column, ""}
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
			if row == Size {
				break
			}
			lastRow := row
			for ; lastRow < Size && xword[lastRow][column] == '-'; lastRow++ {
				//
			}
			if lastRow > row+1 {
				slot := Slot{false, row, column, lastRow - row, ""}
				slots = append(slots, slot)
				row = lastRow
			} else {
				row++
			}
		}
	}

	return slots
}

func conflict(x, y Slot) bool {
	// There's no conflict.
	if x.across == y.across {
		return false
	}

	if !x.across {
		return conflict(y, x)
	}

	if x.row >= y.row &&
		x.row < y.row+y.length &&
		x.column <= y.column &&
		x.column+x.length > y.column {
		if x.word[y.column-x.column] != y.word[x.row-y.row] {
			return true
		}
	}
	return false
}

func checkNewSlot(slot Slot, slots []Slot) bool {
	for _, o := range slots {
		if slot.across == o.across {
			continue
		}
		if conflict(slot, o) {
			return false
		}
	}
	return true
}

func recurse(occupied []Slot, open []Slot, words []string) ([]Slot, bool) {
	// TODO: Test mismatches.
	if len(words) == 0 {
		return occupied, true
	}
	// Don't reorder the slots. Shuffle the strings!
	for i, w := range words {

		if len(w) != open[0].length {
			continue
		}

		// TODO: Have the function return a new one.
		fill := open[0]
		fill.word = w

		if !checkNewSlot(fill, occupied) {
			continue
		}

		remainingWords := append([]string{}, words[:i]...)
		remainingWords = append(remainingWords, words[i+1:]...)
		remainingSlots := append([]Slot{}, occupied...)
		remainingSlots = append(remainingSlots, fill)

		if attempt, ok := recurse(remainingSlots, open[1:], remainingWords); ok {
			return attempt, true
		}

		// Return if one works.
	}

	return nil, false
}

func crosswordPuzzle(puzzle []string, words []string) []string {
	slots := findSlots(puzzle)
	result, _ := recurse(nil, slots, words)
	return render(result)
}

func TestSamples(t *testing.T) {
	type TestCase struct {
		puzzle []string
		words  []string
	}
	table := []TestCase{
		{
			[]string{
				"+-++++++++",
				"+-++++++++",
				"+-++++++++",
				"+-----++++",
				"+-+++-++++",
				"+-+++-++++",
				"+++++-++++",
				"++------++",
				"+++++-++++",
				"+++++-++++",
			},
			strings.Split("LONDON;DELHI;ICELAND;ANKARA", ";"),
		},
		{
			[]string{
				"+-++++++++",
				"+-++++++++",
				"+-------++",
				"+-++++++++",
				"+-++++++++",
				"+------+++",
				"+-+++-++++",
				"+++++-++++",
				"+++++-++++",
				"++++++++++",
			},
			strings.Split("AGRA;NORWAY;ENGLAND;GWALIOR", ";"),
		},
		{
			[]string{
				"++++++-+++",
				"++------++",
				"++++++-+++",
				"++++++-+++",
				"+++------+",
				"++++++-+-+",
				"++++++-+-+",
				"++++++++-+",
				"++++++++-+",
				"++++++++-+",
			},
			strings.Split("ICELAND;MEXICO;PANAMA;ALMATY", ";"),
		},
	}

	for _, row := range table {
		words := append([]string{}, row.words...)
		rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })
		answer := crosswordPuzzle(row.puzzle, row.words)
		t.Logf("%v", strings.Join(answer, "\n"))
	}
}
