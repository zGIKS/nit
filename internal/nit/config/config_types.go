package config

type ClipboardMode string

const (
	ClipboardOnlyCopy ClipboardMode = "only_copy"
	ClipboardAuto     ClipboardMode = "auto"
	ClipboardOSC52    ClipboardMode = "osc52"
	ClipboardSystem   ClipboardMode = "system"
	ClipboardInternal ClipboardMode = "internal"
)

type KeyBinding struct {
	Keys []string `toml:"keys"`
}

type KeyConfig struct {
	Quit         KeyBinding            `toml:"quit"`
	TogglePanel  KeyBinding            `toml:"toggle_panel"`
	FocusCommand KeyBinding            `toml:"focus_command"`
	Down         KeyBinding            `toml:"down"`
	Up           KeyBinding            `toml:"up"`
	ToggleOne    KeyBinding            `toml:"toggle_one"`
	StageAll     KeyBinding            `toml:"stage_all"`
	UnstageAll   KeyBinding            `toml:"unstage_all"`
	Fetch        KeyBinding            `toml:"fetch"`
	Push         KeyBinding            `toml:"push"`
	MenuRight    KeyBinding            `toml:"menu_right"`
	MenuLeft     KeyBinding            `toml:"menu_left"`
	CommitEditor CommitEditorKeyConfig `toml:"commit_editor"`
}

type CommitEditorKeyConfig struct {
	Submit    KeyBinding `toml:"submit"`
	Cancel    KeyBinding `toml:"cancel"`
	Copy      KeyBinding `toml:"copy"`
	Cut       KeyBinding `toml:"cut"`
	Paste     KeyBinding `toml:"paste"`
	SelectAll KeyBinding `toml:"select_all"`
	Backspace KeyBinding `toml:"backspace"`
	Delete    KeyBinding `toml:"delete"`
	Left      KeyBinding `toml:"left"`
	Right     KeyBinding `toml:"right"`
	Home      KeyBinding `toml:"home"`
	End       KeyBinding `toml:"end"`
}

type ClipboardConfig struct {
	Mode     ClipboardMode `toml:"mode"`
	CopyCmd  string        `toml:"copy_cmd"`
	PasteCmd string        `toml:"paste_cmd"`
}

type UIConfig struct {
	RepoLabel                string `toml:"repo_label"`
	BranchLabel              string `toml:"branch_label"`
	RepoBranchSeparator      string `toml:"repo_branch_separator"`
	FetchLabel               string `toml:"fetch_label"`
	MenuLabel                string `toml:"menu_label"`
	MenuChevron              string `toml:"menu_chevron"`
	MenuSelectionIndicator   string `toml:"menu_selection_indicator"`
	BranchSourceSelectedMark string `toml:"branch_source_selected_mark"`
	BranchCreateTitle        string `toml:"branch_create_title"`
	BranchCreateEnterHint    string `toml:"branch_create_enter_hint"`
	BranchCreatePushHint     string `toml:"branch_create_push_hint"`
	BranchCreateNameLabel    string `toml:"branch_create_name_label"`
	BranchCreateSourceLabel  string `toml:"branch_create_source_label"`
}

type FileConfig struct {
	Clipboard ClipboardConfig `toml:"clipboard"`
	Keys      KeyConfig       `toml:"keys"`
	UI        UIConfig        `toml:"ui"`
}

type AppConfig struct {
	ConfigFile       string
	Clipboard        ClipboardConfig
	Keys             KeyConfig
	CommitEditorKeys CommitEditorKeyConfig
	UI               UIConfig
}
