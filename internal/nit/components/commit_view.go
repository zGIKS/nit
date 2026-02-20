package components

func CommitView(width, height int, message string, active bool, submitKey string) string {
	if message == "" {
		message = "Message (" + submitKey + " to commit)"
	}
	lines := []string{
		"[ " + message + " ]",
		"[ Commit ]",
	}
	return BoxView("Changes - Commit", width, height, lines, -1, 0, active, submitKey+" commit")
}
