package state

func insertTextAtCursor(value *string, cursor *int, text string) {
	if value == nil || cursor == nil || text == "" {
		return
	}
	r := []rune(*value)
	if *cursor < 0 {
		*cursor = 0
	}
	if *cursor > len(r) {
		*cursor = len(r)
	}
	insert := []rune(text)
	out := make([]rune, 0, len(r)+len(insert))
	out = append(out, r[:*cursor]...)
	out = append(out, insert...)
	out = append(out, r[*cursor:]...)
	*value = string(out)
	*cursor += len(insert)
}

func backspaceAtCursor(value *string, cursor *int) {
	if value == nil || cursor == nil {
		return
	}
	r := []rune(*value)
	if len(r) == 0 || *cursor <= 0 {
		return
	}
	if *cursor > len(r) {
		*cursor = len(r)
	}
	out := make([]rune, 0, len(r)-1)
	out = append(out, r[:*cursor-1]...)
	out = append(out, r[*cursor:]...)
	*value = string(out)
	*cursor--
}

func deleteAtCursor(value *string, cursor *int) {
	if value == nil || cursor == nil {
		return
	}
	r := []rune(*value)
	if len(r) == 0 {
		return
	}
	if *cursor < 0 {
		*cursor = 0
	}
	if *cursor >= len(r) {
		return
	}
	out := make([]rune, 0, len(r)-1)
	out = append(out, r[:*cursor]...)
	out = append(out, r[*cursor+1:]...)
	*value = string(out)
}

func moveCursorLeft(cursor *int) {
	if cursor != nil && *cursor > 0 {
		*cursor--
	}
}

func moveCursorRight(value string, cursor *int) {
	if cursor == nil {
		return
	}
	if *cursor < len([]rune(value)) {
		*cursor++
	}
}

func moveCursorHome(cursor *int) {
	if cursor != nil {
		*cursor = 0
	}
}

func moveCursorEnd(value string, cursor *int) {
	if cursor != nil {
		*cursor = len([]rune(value))
	}
}

func backspaceWordAtCursor(value *string, cursor *int) {
	if value == nil || cursor == nil {
		return
	}
	r := []rune(*value)
	if len(r) == 0 || *cursor <= 0 {
		return
	}
	if *cursor > len(r) {
		*cursor = len(r)
	}

	i := *cursor - 1
	for i >= 0 && (r[i] == ' ' || r[i] == '\t' || r[i] == '\n') {
		i--
	}
	for i >= 0 && r[i] != ' ' && r[i] != '\t' && r[i] != '\n' {
		i--
	}
	startIdx := i + 1
	if startIdx < 0 {
		startIdx = 0
	}

	out := make([]rune, 0, len(r)-(*cursor-startIdx))
	out = append(out, r[:startIdx]...)
	out = append(out, r[*cursor:]...)
	*value = string(out)
	*cursor = startIdx
}
