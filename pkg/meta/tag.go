package meta

import (
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type Tag []string

// Contains is a wrapper for slices.Contains that returns true if
// the value is explicitly in the tag
func (t Tag) Contains(value string) bool {
	return slices.Contains(t, value)
}

// False only returns true if the tag's first value is explicitly
// set to false or -
func (t Tag) False() bool {
	return slices.Contains([]string{"false", "-"}, t[0])
}

// True returns true if the first value isn't false or "-"
func (t Tag) True() bool {
	return !slices.Contains([]string{"false", "-"}, t[0])
}

// Index is a wrapper for slices.Index that returns the index
// of value in the tag, or -1 if not present
func (t Tag) Index(value string) int {
	return slices.Index(t, value)
}

type Tags map[string]Tag

func ToTags(s string) Tags {
	pattern := regexp.MustCompile(`(?m)(?P<key>\w+):\"(?P<value>[^"]+)\"`)
	matches := pattern.FindAllStringSubmatch(s, -1)
	var tkv = map[string]Tag{}
	for _, match := range matches {
		tkv[match[1]] = strings.Split(match[2], ",")
	}
	return tkv
}

// Value returns the parsed Tag for the given key, or nil if missing
func (t Tags) Value(key string) Tag {
	values, _ := t[key]
	return values
}

// False only returns true if the tags exists and the first value is false or "-"
func (t Tags) False(key string) bool {
	if tag, ok := t[key]; ok && tag != nil {
		return tag.False()
	}
	return false
}

// True returns true if the tag exists and the first value is not false or "-"
func (t Tags) True(key string) bool {
	if tag, ok := t[key]; ok && tag != nil {
		return tag.True()
	}
	return false
}

// Exists returns true if the tag with key exists
func (t Tags) Exists(key string) bool {
	tag, ok := t[key]
	return ok && tag != nil
}

func (t Tags) Contains(key, value string) bool {
	if tag, ok := t[key]; ok && tag != nil {
		return tag.Contains(value)
	}
	return false
}

// Tag returns the tag for key, or nil if it is missing
// note: it is the same as Value
func (t Tags) Tag(key string) Tag {
	tag, _ := t[key]
	return tag
}
