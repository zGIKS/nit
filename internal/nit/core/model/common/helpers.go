package common

import g "nit/internal/nit/git"

func SameChanges(a, b []g.ChangeEntry) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Raw != b[i].Raw {
			return false
		}
	}
	return true
}
