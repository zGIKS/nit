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
