package git

import "testing"

func TestPrettifyGraphLine(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "commit line with graph prefix and hash",
			in:   "| * 8fd9242 message / keep slash in message",
			want: "│ ● 8fd9242 message / keep slash in message",
		},
		{
			name: "continuation graph line without hash",
			in:   "|\\  ",
			want: "│╲  ",
		},
		{
			name: "underscore bridge in graph",
			in:   "| *- not expected but _ bridge",
			want: "│ ●- not expected but _ bridge",
		},
		{
			name: "line without graph prefix stays unchanged",
			in:   "feature/test path and * text",
			want: "feature/test path and * text",
		},
		{
			name: "graph only with diagonal and spaces",
			in:   "| |/ ",
			want: "│ │╱ ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prettifyGraphLine(tt.in)
			if got != tt.want {
				t.Fatalf("prettifyGraphLine(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

