package add

import "testing"

func TestAdd(t *testing.T) {
	a := 2

	b := 3

	exp := 5
	res := add(a, b)
	if exp != res {
		t.Errorf("Expected %d, got %d.", exp, res)
	}
}
