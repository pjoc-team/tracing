package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func ToStr(dval interface{}) string {
	if dval == nil {
		return ""
	}
	switch vv := dval.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", vv)
	case float32, float64:
		return fmt.Sprintf("%g", vv)
	case string:
		return vv
	default:
		return ""
	}
}

func MapToStr(m map[string]interface{}, filter ...string) string {
	var buf strings.Builder
	buf.WriteByte('[')
	for k, v := range m {
		stop := false
		for _, vf := range filter {
			if k == vf {
				stop = true
			}
		}
		if stop {
			continue
		}
		if buf.Len() > 1 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(ToStr(v))
	}
	buf.WriteByte(']')
	v := buf.String()
	if v == "[]" {
		return ""
	}
	return v
}

func GetJsonBytes(i interface{}) []byte {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(i)
	if err != nil {
		return []byte("{}")
	}
	return buffer.Bytes()
}
