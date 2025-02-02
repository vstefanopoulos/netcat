package netcat

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type server struct {
	name          string
	port          string
	clients       map[net.Conn]client
	messages      []message
	welcome       []byte
	mutex         sync.RWMutex
	broadcastChan chan struct{}
	updateUser    chan struct{}
	quit          chan struct{}
	logChan       chan message
}

type client struct {
	name string
	conn net.Conn
	addr net.Addr
}

type message struct {
	senderName string
	senderConn net.Conn
	timeStamp  string
	message    string
}

// Starts the listener, message log, shutdown handler and instantiate wait groups
func (s *server) startServer() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.messageLog()
	}()

	go s.shutDown()

	go s.listen(listener)

	wg.Wait() // wait for log() to finish
	listener.Close()
	s.quit <- struct{}{} // send message to Init()
	return fmt.Errorf("Server shut down")
}

// Listens to incoming connections in for loop. The loop is broken only if the listener is closed
func (s *server) listen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ne, ok := err.(*net.OpError); ok && ne.Op == "accept" {
				break
			}
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go func(c net.Conn) {
			s.handleConnection(c)
		}(conn)
	}
}

// Signals exit msg to log and closes log chan
// Shut down signal spawns from GUI's handling of ^C and results as a quit error in Init()
func (s *server) shutDown() {
	<-s.quit
	msg := message{senderName: "system", message: "Server Closed", timeStamp: "GetOut"}
	s.logChan <- msg
	<-s.logChan
	close(s.logChan)
}

// Initialise server and GUI and oversee graceful shut down with ^C
func Init(name string, port string) {
	s := newServer(name, port)
	go func() {
		if err := s.startServer(); err != nil {
			fmt.Printf("Server: %v\n", err)
		}
	}()

	if err := SetupGui(s); err != nil {
		// When ^C is pressed signals through s.quit to shutDown()
		// Then waits signal back from the channel from StartServer() before returning
		if err.Error() == "quit" {
			s.quit <- struct{}{}
			<-s.quit
			return
		} else {
			log.Printf("Error starting GUI: %v\n", err)
			<-s.quit
			return
		}
	}
	return
}

func newServer(name, port string) *server {
	var welcome strings.Builder
	welcome.WriteString("Welcome to TCP-Chat!\n")
	welcome.WriteString("         _nnnn_\n")
	welcome.WriteString("        dGGGGMMb\n")
	welcome.WriteString("       @p~qp~~qMb\n")
	welcome.WriteString("       M|@||@) M|\n")
	welcome.WriteString("       @,----.JM|\n")
	welcome.WriteString("      JS^\\__/  qKL\n")
	welcome.WriteString("     dZP        qKRb\n")
	welcome.WriteString("    dZP          qKKb\n")
	welcome.WriteString("   fZP            SMMb\n")
	welcome.WriteString("   HZM            MMMM\n")
	welcome.WriteString("   FqM            MMMM\n")
	welcome.WriteString(" __| \".        |\\dS\"qML\n")
	welcome.WriteString(" |    `.       | ' \\Zq\n")
	welcome.WriteString("_)      \\.___.,|     .'\n")
	welcome.WriteString("\\____   )MMMMMP|   .'\n")
	welcome.WriteString("     `-'       `--'\n\n")
	welcome.WriteString("Username rules:\n")
	welcome.WriteString("- Max 12 characters\n")
	welcome.WriteString("- No spaces\n")
	welcome.WriteString("- Unique name\n\n")
	welcome.WriteString("[ENTER YOUR NAME]: ")

	return &server{
		name:          name,
		port:          port,
		clients:       make(map[net.Conn]client),
		welcome:       []byte(welcome.String()),
		broadcastChan: make(chan struct{}, 10),
		updateUser:    make(chan struct{}, 10),
		quit:          make(chan struct{}, 1),
		logChan:       make(chan message, 10),
	}
}
