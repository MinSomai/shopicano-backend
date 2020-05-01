package utils

import "strings"

func IsImage(name string) bool {
	return strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpg") ||
		strings.HasSuffix(name, ".jpeg")
}
