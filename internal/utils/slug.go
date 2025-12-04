package utils

import (
    "regexp"
    "strings"
)

var slugRegex = regexp.MustCompile(`[^a-z0-9\-]+`)

func ToSlug(name string) string {
    s := strings.ToLower(name)
    s = strings.ReplaceAll(s, " ", "-")
    s = slugRegex.ReplaceAllString(s, "")
    return s
}
