package models

import "fmt"

// StringDiff represents a difference on an integer value
type StringDiff struct {
	was   *string
	isNow *string
}

// Empty check if there's no change on the value
func (diff *StringDiff) Empty() bool {
	if diff == nil || diff.was == diff.isNow {
		return true
	}
	if diff.was == nil || diff.isNow == nil {
		return false
	}
	return *diff.was == *diff.isNow
}

func (diff *StringDiff) change(was *string, isNow *string) {
	diff.was = was
	diff.isNow = isNow
}

func (diff *StringDiff) String() string {
	if diff.Empty() {
		return ""
	}
	if diff.was == nil && diff.isNow != nil {
		return fmt.Sprintf("was: <nil> and is now: \"%s\"", *diff.isNow)
	}
	if diff.was != nil && diff.isNow == nil {
		return fmt.Sprintf("was: \"%s\" and is now: <nil>", *diff.was)
	}
	return fmt.Sprintf("was: \"%s\"  and is now: \"%s\"", *diff.was, *diff.isNow)
}
