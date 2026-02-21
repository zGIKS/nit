package git

type ChangeEntry struct {
	X       byte
	Y       byte
	Path    string
	Raw     string
	Staged  bool
	Changed bool
}
