package logger

type FormatType int32

const (
	FormatText FormatType = 0
	FormatJson FormatType = 1
	// Deprecated: 请使用time.RFC3339Nano
	TimeFormatRFC3339Milli = "2006-01-02T15:04:05Z07:00.000"
)

func (format FormatType) String() string {
	switch format {
	case FormatText:
		return "TEXT"
	case FormatJson:
		return "JSON"
	default:
		return "UNKNOWN"
	}
}
