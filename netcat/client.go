package netcat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// Creates new client struct instace, handles log in and messaging for new connection
func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()
	c := client{conn: conn}

	s.mutex.RLock()
	if len(s.clients) == 10 {
		c.conn.Write([]byte("Server full. Please try again later\n"))
		s.mutex.RUnlock()
		return
	}
	s.mutex.RUnlock()

	c.conn.Write(s.welcome)

	var loggedIn bool

	if err := c.validateUName(s); err != nil {
		c.conn.Write([]byte(err.Error()))
		return
	}

	s.sendHistory(c)
	s.addUser(c)
	loggedIn = true

	for {
		msg, err := c.readMsg(loggedIn)
		if err != nil {
			s.logOutUser(c)
			return
		}
		if msg.message == "" {
			c.prompt()
		} else {
			msg.senderConn = conn

			s.mutex.Lock()
			s.broadcastMessage(msg)
			c.prompt()
			s.mutex.Unlock()
		}
	}
}

// Returns message struct instance with input from terminal if client is logged in
// If client is loging in it updates the client struct with given name
func (c *client) readMsg(loggedIn bool) (message, error) {
	var msg message
	reader := bufio.NewReader(c.conn)
	input, err := reader.ReadString('\n')
	if err != nil {
		return msg, fmt.Errorf("Invalid input")
	}
	input = strings.TrimSpace(input)

	if input == "" {
		return msg, nil
	}

	if !loggedIn {
		// Return raw input for validation
		return message{message: input}, nil
	} else {
		msg = message{senderName: c.name,
			timeStamp:  formatTimestamp(),
			message:    input,
			senderConn: c.conn}
	}
	return msg, nil
}

// Writes to connection with a set deadline
func (c *client) write(msg string) {
	var mx sync.Mutex
	mx.Lock()
	defer mx.Unlock()
	c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		return
	}
}

// Checks for non and duplicate names, spaces in name and for length over 12 chars and updates client name.
// Allows only five attempts
func (c *client) validateUName(s *server) error {
	var count int
	for {
		count++
		if count > 5 {
			return fmt.Errorf("\nToo many invalid attempts. You are disconnected\n")
		}

		nameMsg, err := c.readMsg(false)
		if err != nil || nameMsg.message == "" {
			c.conn.Write([]byte("Invalid input. Try again: "))
			continue
		}

		if len(nameMsg.message) > 12 {
			c.conn.Write([]byte("Name too long (max 12 chars). Try again: "))
			continue
		}

		if strings.Contains(nameMsg.message, " ") {
			c.conn.Write([]byte("No spaces allowed in username. Try again: "))
			continue
		}

		// Check uniqueness
		s.mutex.RLock()
		unique := true
		for _, cl := range s.clients {
			if cl.name == nameMsg.message {
				unique = false
				break
			}
		}
		s.mutex.RUnlock()

		if !unique {
			c.conn.Write([]byte("Name already taken. Try again: "))
			continue
		}

		c.name = nameMsg.message
		return nil
	}
}
