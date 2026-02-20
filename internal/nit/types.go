package nit

type panelMode int

type uiMode int

type paneFocus int

const (
	panelGraph panelMode = iota
	panelOutput
)

const (
	uiBrowse uiMode = iota
	uiMenu
	uiPrompt
)

const (
	focusGraph paneFocus = iota
	focusChanges
)

type action struct {
	label string
	kind  string
}

type changeEntry struct {
	x       byte
	y       byte
	path    string
	raw     string
	staged  bool
	changed bool
}

type cmdResultMsg struct {
	title          string
	output         []string
	err            error
	switchToOutput bool
}

type keyBinding struct {
	Keys []string `json:"keys"`
}

type keyConfig struct {
	Quit            keyBinding `json:"quit"`
	OpenMenu        keyBinding `json:"open_menu"`
	TogglePanel     keyBinding `json:"toggle_panel"`
	ShowOutput      keyBinding `json:"show_output"`
	Reload          keyBinding `json:"reload"`
	Down            keyBinding `json:"down"`
	Up              keyBinding `json:"up"`
	PageDown        keyBinding `json:"page_down"`
	PageUp          keyBinding `json:"page_up"`
	Home            keyBinding `json:"home"`
	End             keyBinding `json:"end"`
	StageSelected   keyBinding `json:"stage_selected"`
	UnstageSelected keyBinding `json:"unstage_selected"`
	StageAll        keyBinding `json:"stage_all"`
	UnstageAll      keyBinding `json:"unstage_all"`
	MenuDown        keyBinding `json:"menu_down"`
	MenuUp          keyBinding `json:"menu_up"`
	MenuSelect      keyBinding `json:"menu_select"`
	MenuClose       keyBinding `json:"menu_close"`
	PromptSubmit    keyBinding `json:"prompt_submit"`
	PromptCancel    keyBinding `json:"prompt_cancel"`
	PromptBackspace keyBinding `json:"prompt_backspace"`
}

type model struct {
	ui    uiMode
	panel panelMode
	focus paneFocus

	graphLines    []string
	outputLines   []string
	changeEntries []changeEntry
	changeLines   []string
	lines         []string

	cursor int
	offset int
	width  int
	height int

	changesCursor int
	changesOffset int

	err    string
	status string

	menuItems  []action
	menuCursor int

	promptTitle       string
	promptPlaceholder string
	promptValue       string
	promptKind        string

	keys keyConfig
}
