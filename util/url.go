package util

import (
	"net/url"
	"strings"
)

func GetPath(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	return u.Path
}

func GetPreAndFileName(url string) string {
	if url == "" {
		return url
	}
	paths := strings.Split(url, "/")
	l := len(paths)
	if l >= 2 {
		var buf strings.Builder
		head := paths[l-2]
		if head != "" {
			buf.WriteString(head)
			buf.WriteString("/")
		}
		buf.WriteString(paths[l-1])
		return buf.String()
	}
	return url
}
