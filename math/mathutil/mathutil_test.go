package mathutil

import (
	"testing"
)

func TestPercent(t *testing.T) {
	cases := []struct {
		numerator   int64
		denominator int64
		normalize   int64
		formatted   string
	}{
		{0, 0, 0, "0.00%"},
		{0, 1, 0, "0.00%"},
		{0, -1, 0, "0.00%"},
		{1, 1, 10000, "100.00%"},
		{1, -1, -10000, "-100.00%"},
		{-1, 1, -10000, "-100.00%"},
		{-1, -1, 10000, "100.00%"},
		{10, 100000, 1, "0.01%"},
		{10, 10000, 10, "0.10%"},
		{10, 100, 1000, "10.00%"},
		{1, 100000, 0, "0.00%"},
	}
	for i, c := range cases {
		p := NewPercent(c.numerator, c.denominator)
		normalize := p.NormalizeWithBase(10000)
		formatted := p.String()
		if normalize != c.normalize {
			t.Errorf("case %02d: normalize want %d, but got %d", i, c.normalize, normalize)
		}
		if formatted != c.formatted {
			t.Errorf("case %02d: formatted want %s, but got %s", i, c.formatted, formatted)
		}
	}
}

func TestComb(t *testing.T) {
	const m = 10
	var cases = map[int]int64{
		0:  1,
		1:  10,
		2:  45,
		3:  120,
		4:  210,
		5:  252,
		6:  210,
		7:  120,
		8:  45,
		9:  10,
		10: 1,
	}
	for n, sum := range cases {
		got := Comb(m, n)
		if got != sum {
			t.Errorf("C(%d,%d) want %d, but got %d", m, n, sum, got)
		}
		result := CombSet(m, n)
		if len(result) != int(sum) {
			t.Errorf("C(%d,%d) want %d, but got len(result)=%d", m, n, sum, len(result))
		}
	}
}

func TestMultiComb(t *testing.T) {
	nums := []int{3, 2, 1, 1}
	for i := 0; i <= 7; i++ {
		t.Logf("MC(%v, %d): %v", nums, i, MultiCombSet(nums, i))
	}
}
