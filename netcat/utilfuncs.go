package netcat

import (
	"strings"
	"time"
)

// Format timestamp for message
func formatTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// receives message instance and returns formated string
func (msg *message) formatMessage() string {
	if msg.timeStamp != "" && msg.senderName != "" {
		var str strings.Builder
		str.WriteString("[")
		str.WriteString(msg.timeStamp)
		str.WriteString("]")
		str.WriteString("[")
		str.WriteString(msg.senderName)
		str.WriteString("]: ")
		str.WriteString(msg.message)
		return str.String()
	} else {
		return msg.message
	}
}
