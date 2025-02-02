package netcat

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Creates message log file and iterates through s.logChan until the channel is closed
func (s *server) messageLog() {
	var exit bool

	logName := time.Now().Format("2006-01-02_15-04-05") + ".log"
	f, err := os.Create(logName)
	if err != nil {
		log.Printf("Couldn't create log file: %v", err)
		return
	}

	defer func() {
		f.Close()
		fmt.Println("Log: Closing log file")
	}()

	s.logChan <- message{senderName: "System", message: "Server Started..."}
	for msg := range s.logChan {

		if msg.timeStamp == "GetOut" {
			exit = true
			msg.timeStamp = formatTimestamp()
		}
		if msg.timeStamp == "" {
			msg.timeStamp = formatTimestamp()
		}

		if msg.senderName == "" {
			msg.senderName = "System"
		}

		if f != nil {
			f.WriteString(msg.formatMessage())
			f.WriteString("\n")
			f.Sync() // Flush to disk
		}

		if exit {
			s.logChan <- message{}
			break
		}
	}
}
