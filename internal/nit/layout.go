package nit

func (m *model) clamp() {
	if m.focus == focusGraph {
		if len(m.graphLines) == 0 {
			m.cursor = 0
			m.offset = 0
			return
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
		if m.cursor >= len(m.graphLines) {
			m.cursor = len(m.graphLines) - 1
		}
		page := m.graphPageSize()
		if m.cursor < m.offset {
			m.offset = m.cursor
		}
		if m.cursor >= m.offset+page {
			m.offset = m.cursor - page + 1
		}
		maxOffset := max(0, len(m.graphLines)-page)
		if m.offset > maxOffset {
			m.offset = maxOffset
		}
		if m.offset < 0 {
			m.offset = 0
		}
		return
	}

	if len(m.changeLines) == 0 {
		m.changesCursor = 0
		m.changesOffset = 0
		return
	}
	if m.changesCursor < 0 {
		m.changesCursor = 0
	}
	if m.changesCursor >= len(m.changeLines) {
		m.changesCursor = len(m.changeLines) - 1
	}
	page := m.changesPageSize()
	if m.changesCursor < m.changesOffset {
		m.changesOffset = m.changesCursor
	}
	if m.changesCursor >= m.changesOffset+page {
		m.changesOffset = m.changesCursor - page + 1
	}
	maxOffset := max(0, len(m.changeLines)-page)
	if m.changesOffset > maxOffset {
		m.changesOffset = maxOffset
	}
	if m.changesOffset < 0 {
		m.changesOffset = 0
	}
}

func (m model) bodyHeight() int {
	h := m.height
	if h < 6 {
		return 6
	}
	return h
}

func (m model) graphPaneHeight() int {
	h := m.bodyHeight()
	gh := (h * 45) / 100
	if gh < 4 {
		gh = 4
	}
	if gh > h-4 {
		gh = h - 4
	}
	return gh
}

func (m model) changesPaneHeight() int {
	h := m.bodyHeight()
	ch := h - m.graphPaneHeight()
	if ch < 4 {
		return 4
	}
	return ch
}

func (m model) graphPageSize() int {
	h := m.graphPaneHeight() - 2
	if h < 1 {
		return 1
	}
	return h
}

func (m model) changesPageSize() int {
	h := m.changesPaneHeight() - 2
	if h < 1 {
		return 1
	}
	return h
}
