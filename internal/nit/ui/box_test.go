package ui

import "testing"

func TestDisplayWidthIgnoresANSI(t *testing.T) {
	s := "\x1b[4m[f] fetch\x1b[24m"
	if got, want := displayWidth(s), 9; got != want {
		t.Fatalf("displayWidth() = %d, want %d", got, want)
	}
}

func TestFitTextKeepsVisibleWidthWithANSI(t *testing.T) {
	s := "\x1b[4mabc\x1b[24m"
	got := fitText(s, 5, ' ')
	if w := displayWidth(got); w != 5 {
		t.Fatalf("fitText visible width = %d, want 5; got %q", w, got)
	}
}

