package lib

import (
	"path"
	"strings"
)

// FileExt returns the lowercase extension of name without the leading dot.
// Examples:
//
//	"photo.JPG" -> "jpg"
//	"archive.tar.gz" -> "gz"
//	"noext" -> ""
func FileExt(name string) string {
	if name == "" {
		return ""
	}
	return strings.ToLower(strings.TrimPrefix(path.Ext(name), "."))
}
