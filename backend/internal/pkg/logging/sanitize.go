package logging

import "strings"

// SanitizeString removes line breaks from values that may originate from user input (CWE-117).
func SanitizeString(value string) string {
	value = strings.ReplaceAll(value, "\n", "")
	value = strings.ReplaceAll(value, "\r", "")
	return value
}

func sanitizeLogValue(value any) any {
	switch v := value.(type) {
	case string:
		return SanitizeString(v)
	case error:
		if v == nil {
			return nil
		}
		return SanitizeString(v.Error())
	default:
		return value
	}
}

func sanitizeLogArgs(args ...any) []any {
	if len(args) == 0 {
		return args
	}
	out := make([]any, len(args))
	for i, arg := range args {
		if i%2 == 1 {
			out[i] = sanitizeLogValue(arg)
			continue
		}
		out[i] = arg
	}
	return out
}

func sanitizeLogMessage(msg string) string {
	return SanitizeString(msg)
}
