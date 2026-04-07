package state

func firstSelectableIndex(items []DropdownMenuItem) int {
	for i, item := range items {
		if !item.Separator {
			return i
		}
	}
	return -1
}

func nextSelectableIndex(items []DropdownMenuItem, start, delta int) int {
	if len(items) == 0 || delta == 0 {
		return start
	}
	if start < 0 || start >= len(items) {
		start = firstSelectableIndex(items)
	}
	if start < 0 {
		return -1
	}
	i := start
	for {
		i += delta
		if i < 0 {
			i = len(items) - 1
		}
		if i >= len(items) {
			i = 0
		}
		if !items[i].Separator {
			return i
		}
		if i == start {
			return start
		}
	}
}

func clampScrollSelection(items []DropdownMenuItem, hover *int, offset *int, page int) {
	if len(items) == 0 {
		*hover = -1
		*offset = 0
		return
	}
	if page < 1 {
		page = 1
	}
	if *hover < 0 || *hover >= len(items) || items[*hover].Separator {
		*hover = firstSelectableIndex(items)
	}
	if *hover < 0 {
		*offset = 0
		return
	}
	if *hover < *offset {
		*offset = *hover
	}
	if *hover >= *offset+page {
		*offset = *hover - page + 1
	}
	maxOffset := max(0, len(items)-page)
	if *offset > maxOffset {
		*offset = maxOffset
	}
	if *offset < 0 {
		*offset = 0
	}
}

func menuPageSizeForRectHeight(h int) int {
	page := h - 2
	if page < 1 {
		return 1
	}
	return page
}
