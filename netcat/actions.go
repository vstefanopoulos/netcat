package netcat

import (
	"log"
	"net"
	"strings"
	"time"
)

// Broadcasts a raw string to all online users and appends it on server []messages
// Lock the mutex before calling and unlock after
func (s *server) broadcastMessage(msg message) {
	if msg.message == "" {
		return
	}
	s.messages = append(s.messages, msg)

	go func() { s.logChan <- msg }()

	for _, receiver := range s.clients {
		if receiver.conn != msg.senderConn {
			receiver.write("\n")
			receiver.write(msg.formatMessage())
		}
	}

	s.promptAllExcept(msg.senderConn)

	select {
	case s.broadcastChan <- struct{}{}:
	case <-time.After(time.Second):
		log.Println("Update user event skipped due to full channel")
	}
}

// Sends all the messages to new connection and then prompts.
// Concurrency safe
func (s *server) sendHistory(c client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, msg := range s.messages {
		c.write(msg.formatMessage())
		c.write("\n")
	}
	c.prompt()
}

// Broadcast message to all clients and updates server clients map.
// Prompts all clients except message sender conn.
// Concurrency safe
func (s *server) addUser(c client) {
	var msg strings.Builder
	msg.WriteString(c.name)
	msg.WriteString(" has joined our chat...")

	s.mutex.Lock()

	m := message{senderConn: c.conn,
		message: msg.String()}
	s.broadcastMessage(m)
	s.clients[c.conn] = c

	s.mutex.Unlock()

	select {
	case s.updateUser <- struct{}{}:
	case <-time.After(time.Second):
		log.Println("Update user event skipped due to full channel")
	}
}

// Removes user from server clients and broadcasts message.
// Concurrency safe
func (s *server) logOutUser(c client) {
	var msg strings.Builder
	msg.WriteString(c.name)
	msg.WriteString(" has left our chat...")

	s.mutex.Lock()
	delete(s.clients, c.conn)
	m := message{senderConn: c.conn,
		message: msg.String()}

	s.broadcastMessage(m)
	s.mutex.Unlock()

	select {
	case s.updateUser <- struct{}{}:
	case <-time.After(time.Second):
		log.Println("Update user event skipped due to full channel")
	}
}

// Prompt client with timestamp and name
func (c client) prompt() {
	msg := message{senderName: c.name, timeStamp: formatTimestamp()}
	c.write(msg.formatMessage())
}

// Prompt all users except sender adding a new line before the prompt.
// Lock the mutex before calling and unlock after
func (s *server) promptAllExcept(sender net.Conn) {
	for _, receiver := range s.clients {
		if receiver.conn != sender {
			receiver.write("\n")
			receiver.prompt()
		}
	}
}
