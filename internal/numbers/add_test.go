package numbers

import "testing"

func TestAdd(t *testing.T) {
	cases := []struct {
		inA, inB, want int
	}{
		{0, 1, 1},
	}
	for _, c := range cases {
		got := Add(c.inA, c.inB)
		if got != c.want {
			t.Errorf("%q, %q == %q, want %q", c.inA, c.inB, got, c.want)
		}
	}
}
