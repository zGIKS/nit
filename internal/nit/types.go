package nit

type paneFocus int

const (
	focusChanges paneFocus = iota
	focusGraph
)

type changeEntry struct {
	x       byte
	y       byte
	path    string
	raw     string
	staged  bool
	changed bool
}

type changeRow struct {
	text       string
	selectable bool
	section    string
	index      int
}

type keyBinding struct {
	Keys []string `json:"keys"`
}

type keyConfig struct {
	Quit        keyBinding `json:"quit"`
	TogglePanel keyBinding `json:"toggle_panel"`
	Down        keyBinding `json:"down"`
	Up          keyBinding `json:"up"`
	ToggleOne   keyBinding `json:"toggle_one"`
	StageAll    keyBinding `json:"stage_all"`
	UnstageAll  keyBinding `json:"unstage_all"`
}

type model struct {
	focus paneFocus

	graphLines      []string
	changeEntries   []changeEntry
	stagedChanges   []changeEntry
	unstagedChanges []changeEntry
	changeRows      []changeRow
	changeLines     []string

	cursor int
	offset int
	width  int
	height int

	changesCursor int
	changesOffset int

	keys keyConfig
}
