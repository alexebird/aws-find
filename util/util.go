package util

func WithDefault(val *string, defaultVal string) string {
	if val != nil {
		return *val
	} else {
		return defaultVal
	}
}

func WithEmptyStringDefault(val *string) string {
	return WithDefault(val, "")
}
