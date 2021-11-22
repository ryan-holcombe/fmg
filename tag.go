package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	tagSkip     = "-"
	tagOptional = "optional"
	tagName     = "fmgen"
)

var tagRegex = regexp.MustCompile(fmt.Sprintf(`.*%s:"([^\)]+)".*`, tagName))

type tag struct {
	values []string
}

func (t tag) skip() bool {
	for _, v := range t.values {
		if v == tagSkip {
			return true
		}
	}
	return false
}

func (t tag) optional() bool {
	for _, v := range t.values {
		if v == tagOptional {
			return true
		}
	}
	return false
}

func parseTag(allTags string) (tag, bool) {
	rs := tagRegex.FindStringSubmatch(allTags)
	if len(rs) <= 1 {
		return tag{}, false
	}

	tagSplit := strings.Split(rs[1], ",")

	var tagValues []string
	for _, t := range tagSplit {
		tagValues = append(tagValues, strings.TrimSpace(t))
	}

	return tag{tagValues}, true
}
