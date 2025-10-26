package components

import "strings"

// stripANSI removes ANSI escape codes from a string for testing
func StripANSI(s string) string {
	// Simple ANSI stripper - removes common escape sequences
	result := s
	// Remove CSI sequences (most common)
	for strings.Contains(result, "\x1b[") {
		start := strings.Index(result, "\x1b[")
		end := start + 2
		for end < len(result) && !((result[end] >= 'A' && result[end] <= 'Z') || (result[end] >= 'a' && result[end] <= 'z')) {
			end++
		}
		if end < len(result) {
			end++
		}
		result = result[:start] + result[end:]
	}
	return result
}
