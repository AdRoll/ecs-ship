package models

func updateString(old *string, new *string, apply func(string), record func(*string, *string)) {
	if old == nil && new == nil || new == nil {
		return
	}
	record(old, new)
	apply(*new)
}

func updateInt(old *int32, new *int32, apply func(int32), record func(*int32, *int32)) {
	if old == nil && new == nil || new == nil {
		return
	}
	record(old, new)
	apply(*new)
}
