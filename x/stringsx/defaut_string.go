package stringsx

func DefaultString(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
