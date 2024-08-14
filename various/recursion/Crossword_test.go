package recursion

import "testing"

type Word struct {
	word        string
	horizontal  bool
	row, column int
}

func (w *Word) collides(x *Word) bool {
	if w.horizontal != x.horizontal {
		return false
	}

	if w.horizontal {

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
