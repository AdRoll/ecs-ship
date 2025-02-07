package models

import "fmt"

// IntegerDiff represents a difference on an integer value
type IntegerDiff struct {
	was   *int32
	isNow *int32
}


// Empty check if there's no change on the value
func (diff *IntegerDiff) Empty() bool {
	if diff == nil || diff.was == diff.isNow {
		return true
	}
	if diff.was == nil || diff.isNow == nil {
		return false
	}
	return *diff.was == *diff.isNow
}

func (diff *IntegerDiff) change(was *int32, isNow *int32) {
	diff.was = was
	diff.isNow = isNow
}

func (diff *IntegerDiff) String() string {
	if diff.Empty() {
		return ""
	}
	if diff.was == nil && diff.isNow != nil {
		return fmt.Sprintf("was: <nil> and is now: %d", *diff.isNow)
	}
	if diff.was != nil && diff.isNow == nil {
		return fmt.Sprintf("was: %d and now is: <nil>", *diff.was)
	}
	return fmt.Sprintf("was: %d and now is: %d", *diff.was, *diff.isNow)
}
