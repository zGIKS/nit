package state

import "strings"

func (s *AppState) syncBranchCreateSources() {
sources := make([]string, 0, len(s.Branches.Lines))
for _, line := range s.Branches.Lines {
name := strings.TrimSpace(line)
if name == "" || strings.Contains(name, "No local branches") || strings.Contains(name, "Loading branches") || strings.Contains(name, "Not a git repo") {
continue
}
name = strings.TrimPrefix(name, "●")
name = strings.TrimPrefix(name, "*")
name = strings.TrimSpace(name)
if name == "" {
continue
}
sources = append(sources, name)
}
s.BranchCreateSourceList = sources
if s.BranchCreateSource == "" && s.BranchName != "" {
s.BranchCreateSource = s.BranchName
}
if findStringIndex(s.BranchCreateSourceList, s.BranchCreateSource) < 0 {
if len(s.BranchCreateSourceList) > 0 {
s.BranchCreateSource = s.BranchCreateSourceList[0]
} else if s.BranchCreateSource == "" {
s.BranchCreateSource = "-"
}
}
s.ensureBranchCreateSourceVisible()
}

func (s AppState) BranchCreateSourceIndexAt(x, y int) (int, bool) {
if !s.BranchCreateOpen {
return -1, false
}
lx, ly, lw, lh := s.BranchCreateSourceListRect()
if x < lx || x >= lx+lw || y < ly || y >= ly+lh {
return -1, false
}
idx := s.BranchCreateSourceOffset + (y - ly)
if idx < 0 || idx >= len(s.BranchCreateSourceList) {
return -1, false
}
return idx, true
}

func (s *AppState) BranchCreateSelectSourceIndex(idx int) {
if idx < 0 || idx >= len(s.BranchCreateSourceList) {
return
}
s.BranchCreateSource = s.BranchCreateSourceList[idx]
s.ensureBranchCreateSourceVisible()
}

func (s *AppState) BranchCreateMoveSource(delta int) {
if len(s.BranchCreateSourceList) == 0 || delta == 0 {
return
}
cur := 0
for i, name := range s.BranchCreateSourceList {
if name == s.BranchCreateSource {
cur = i
break
}
}
cur += delta
if cur < 0 {
cur = 0
}
if cur >= len(s.BranchCreateSourceList) {
cur = len(s.BranchCreateSourceList) - 1
}
s.BranchCreateSource = s.BranchCreateSourceList[cur]
s.ensureBranchCreateSourceVisible()
}

func (s *AppState) BranchCreateHoverAt(x, y int) {
if idx, ok := s.BranchCreateSourceIndexAt(x, y); ok {
s.BranchCreateHoverIndex = idx
return
}
s.BranchCreateHoverIndex = -1
}
