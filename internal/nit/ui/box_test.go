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

func TestBordersDoNotPanicWhenContentFillsWidth(t *testing.T) {
	for _, innerW := range []int{1, 2, 3, 10} {
		t.Run("top", func(t *testing.T) {
			got := renderTopBorder("Commit", "", innerW, true)
			if displayWidth(got) != innerW+2 {
				t.Fatalf("top border width = %d, want %d; got %q", displayWidth(got), innerW+2, got)
			}
		})
		t.Run("bottom", func(t *testing.T) {
			got := renderBottomBorder("error: something went wrong", innerW, true)
			if displayWidth(got) != innerW+2 {
				t.Fatalf("bottom border width = %d, want %d; got %q", displayWidth(got), innerW+2, got)
			}
		})
	}
}
