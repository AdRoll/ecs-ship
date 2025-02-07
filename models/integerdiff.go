package models

import "fmt"

// IntegerDiff represents a difference on an integer value
type IntegerDiff struct {
	wasNil   bool
	isNowNil bool
	was      int32
	isNow    int32
}

// Empty check if there's no change on the value
func (diff *IntegerDiff) Empty() bool {
	if diff == nil {
		return true
	}
	if diff.wasNil && diff.isNowNil {
		return true
	}
	if diff.wasNil || diff.isNowNil {
		return false
	}
	return diff.was == diff.isNow
}

func (diff *IntegerDiff) Change(was *int32, isNow *int32) {
	if was == nil {
		diff.wasNil = true
	} else {
		diff.wasNil = false
		diff.was = *was
	}
	if isNow == nil {
		diff.isNowNil = true
	} else {
		diff.isNowNil = false
		diff.isNow = *isNow
	}
}

func (diff *IntegerDiff) String() string {
	if diff.Empty() {
		return ""
	}
	wasStr := "<nil>"
	if !diff.wasNil {
		wasStr = fmt.Sprintf("%d", diff.was)
	}
	isNowStr := "<nil>"
	if !diff.isNowNil {
		isNowStr = fmt.Sprintf("%d", diff.isNow)
	}
	return fmt.Sprintf("was: %s and now is: %s", wasStr, isNowStr)
}
