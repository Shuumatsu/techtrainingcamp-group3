package models

import "testing"

func TestInt2Str(t *testing.T) {
	var a, b, c uint64
	a = 13215646
	b = 987463843
	c = 1316563135
	if int2str(a) != "13215646" {
		t.Fatal(int2str(a), "!=", "13215646")
	}
	if int2str(b) != "987463843" {
		t.Fatal(int2str(b), "!=", "987463843")
	}
	if int2str(c) != "1316563135" {
		t.Fatal(int2str(c), "!=", "1316563135")
	}
}
