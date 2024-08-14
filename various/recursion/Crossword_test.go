package recursion

import "testing"

type Word struct {
	word        string
	across  bool
	row, column int
}

type Crossword struct {
	words []Word
}

func (xword *Crossword) String() string {
	return "";
}

func (w *Word) collides(x *Word) bool {
	if w.across != x.across {
		return false
	}

	if w.across {

		if w.row != x.row {
			return false
		}

		if w.column <= x.column {
			return w.column+len(w.word) - 1 >= x.column
		} else {
			return x.column+len(x.word) - 1 >= w.column
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

func (w *Word) crosses(x *Word) bool {
	if w.across == x.across {
		return false;
	}

	if !w.across {
		return x.crosses(w)
	}

	for c := w.column; c < w.column + len(w.word); c++ {
		for r := x.row; r < x.row + len(x.word); r++ {
			if w.word[c] == x.word[r] {
				return true;
			}
		}
	}
	return false;
}

func (w *Word) findCrossings(s string) []*Word {
	var crossings []*Word
	if w.across {
		// Pretend that w starts at (0, 0). The math is easier.
		for c := 0; c < len(w.word); c++ {
			for r := -len(s) + 1; r < 1; r++ {
				if w.word[c] == s[r + len(s) - 1] {
					// Now displace the second word.
					crossings = append(crossings, &Word{s, false, r + w.row, c + w.column})
				}
			}
		}
	} else {
		// Pretend that w starts at (0, 0)
		for r := 0; r < len(w.word); r++ {
			for c := -len(s) + 1; c < 1; c++ {
				if w.word[r] == s[c + len(s) - 1] {
					// TODO: Add tests for this.
					crossings = append(crossings, &Word{s, true, c + w.column, r + w.row})
				}
			}
		}
	}

	return crossings;
}

func TestCrossings(t *testing.T) {
	// For any word with no repeated letters, 
	// the number of crossings must equal the
	// number of letters.
	word := Word{"ALIEN", true, 0, 0}
	crossings := word.findCrossings("ALIEN")
	t.Logf("Found %v", crossings)


	badWord := Word{"SYRINX", true, 0, 0}
	noCrossings := badWord.findCrossings("MOPE")
	if len(noCrossings) > 0 {
		t.Error("Should be 0")
	}
}

func TestCollisions(t *testing.T) {
	type TestCase struct {
		w Word;
		x Word;
		expected bool;
	}
	table := []TestCase{
		{Word{ "GWALIOR", true, 0, 0}, Word{"GWALIOR", true, 0, 0}, true },
		{Word{ "GWALIOR", true, 0, 0}, Word{"GWALIOR", false, 0, 0}, false },
		{Word{ "GWALIOR", true, 0, 0}, Word{"GWALIOR", true, 0, 1}, true },
		{Word{ "GWALIOR", true, 0, 0}, Word{"GWALIOR", true, 0, 6}, true },
		{Word{ "GWALIOR", true, 0, 0}, Word{"GWALIOR", true, 0, 7}, false },
	}
	for i, test := range table {
		if test.w.collides(&test.x) != test.expected {
			t.Errorf("Row %d: %v and %v collision expected %v", i, test.w, test.x, test.expected)
		}

		// Test the commutation!
		if test.x.collides(&test.w) != test.expected {
			t.Errorf("Row %d: %v and %v collision (inverted) expected %v", i, test.w, test.x, test.expected)
		}

		// TODO Now switch all the values for "horizontal".
	}
}
