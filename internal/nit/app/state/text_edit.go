package state

func clearSelectedText(value *string, cursor *int, selectAll *bool) bool {
	if value == nil || cursor == nil || selectAll == nil || !*selectAll {
		return false
	}
	*value = ""
	*cursor = 0
	*selectAll = false
	return true
}

func appendTextInput(value *string, cursor *int, selectAll *bool, text string) {
	if text == "" {
		return
	}
	clearSelectedText(value, cursor, selectAll)
	insertTextAtCursor(value, cursor, text)
}

func backspaceTextInput(value *string, cursor *int, selectAll *bool) {
	if clearSelectedText(value, cursor, selectAll) {
		return
	}
	backspaceAtCursor(value, cursor)
}

func deleteTextInput(value *string, cursor *int, selectAll *bool) {
	if clearSelectedText(value, cursor, selectAll) {
		return
	}
	deleteAtCursor(value, cursor)
}

func moveTextInputCursorLeft(cursor *int, selectAll *bool) {
	if selectAll != nil {
		*selectAll = false
	}
	moveCursorLeft(cursor)
}

func moveTextInputCursorRight(value string, cursor *int, selectAll *bool) {
	if selectAll != nil {
		*selectAll = false
	}
	moveCursorRight(value, cursor)
}

func moveTextInputCursorHome(cursor *int, selectAll *bool) {
	if selectAll != nil {
		*selectAll = false
	}
	moveCursorHome(cursor)
}

func moveTextInputCursorEnd(value string, cursor *int, selectAll *bool) {
	if selectAll != nil {
		*selectAll = false
	}
	moveCursorEnd(value, cursor)
}

func selectAllTextInput(value string, cursor *int, selectAll *bool) {
	if cursor == nil || selectAll == nil {
		return
	}
	if value == "" {
		*selectAll = false
		*cursor = 0
		return
	}
	*selectAll = true
	moveCursorEnd(value, cursor)
}

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
