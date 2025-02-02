package netcat

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// Creates gui and listens on the server channels
func SetupGui(s *server) error {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer gui.Close()

	// Initialise the views layout
	gui.SetManagerFunc(Layout)

	// Set keybindings
	if err := SetupKeybindings(gui, s); err != nil {
		return err
	}

	// Start listening for updates to message or user log
	go StartListening(s, gui)

	// Run the GUI loop
	return gui.MainLoop()
}

// Sets keybindings for <ctrl-C> -> quit and <enter> -> send message from admin
func SetupKeybindings(gui *gocui.Gui, s *server) error {
	// Set the keybinding for Ctrl+C (quit)
	if err := gui.SetKeybinding("input", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return fmt.Errorf("failed to set Ctrl+C keybinding: %w", err)
	}
	// Set the keybinding for Enter (input)
	if err := gui.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return SendMessage(g, v, s)
	}); err != nil {
		return fmt.Errorf("failed to set Enter keybinding: %w", err)
	}

	return nil
}

// Creates layout with "left" view for messages "right" for users and "input" for sending messages from admin
func Layout(gui *gocui.Gui) error {
	width, height := gui.Size()

	// Create the admin panel
	if _, err := gui.SetView("left", 1, 1, width-31, height-5); err != nil && err != gocui.ErrUnknownView {
		return err
	}

	// Set title for the left panel
	if v, err := gui.View("left"); err == nil {
		v.Frame = true
		v.Title = "Admin Tools"
	}

	// Create the channels+users panel
	if _, err := gui.SetView("right", width-30, 1, width-1, height-5); err != nil && err != gocui.ErrUnknownView {
		return err
	}

	// Set title for the right panel
	if v, err := gui.View("right"); err == nil {
		v.Frame = true
		v.Title = "Connected Users"
	}

	// Create the input panel
	if _, err := gui.SetView("input", 1, height-4, width-1, height-1); err != nil && err != gocui.ErrUnknownView {
		return err
	}

	// Set title for the input panel
	if v, err := gui.View("input"); err == nil {
		v.Frame = true
		v.Title = "Input"
		v.Editable = true
		v.Wrap = false
	}

	// Set the input view as the focused view
	_, err := gui.SetCurrentView("input")
	if err != nil {
		return err
	}

	return nil
}

// Broadcasts a message specifically from the admin
func SendMessage(gui *gocui.Gui, v *gocui.View, s *server) error {
	text := v.Buffer()
	v.Clear()

	// Reset the cursor position to the beginning
	v.SetCursor(0, 0)

	if text == "" {
		return nil
	}

	// Remove the newline
	text = text[:len(text)-1]

	// Prepare the message
	senderName := "Admin"
	timeStamp := formatTimestamp()
	msg := message{senderName: senderName, timeStamp: timeStamp, message: text, senderConn: nil}

	// Log the message to the server
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.broadcastMessage(msg)

	return nil
}

// Close the GUI and server
func Quit(gui *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// Start goroutines to listen for updates and refresh the GUI
func StartListening(s *server, g *gocui.Gui) {
	for {
		select {
		case <-s.broadcastChan:
			g.Update(func(g *gocui.Gui) error {
				refreshFunc(g, "left", s)
				return nil
			})
		case <-s.updateUser:
			g.Update(func(g *gocui.Gui) error {
				refreshFunc(g, "right", s)
				return nil
			})
		}
	}
}

// Updates the corresponding window when the message log or the user log updates
func refreshFunc(gui *gocui.Gui, viewName string, s *server) {
	v, err := gui.View(viewName)
	if err != nil {
		log.Println("Error fetching view:", viewName)
		return
	}
	v.Clear()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	switch viewName {
	case "left":
		// Refresh the left panel (Admin Tools, message logs)
		numMessages := len(s.messages)
		if numMessages > 20 {
			numMessages = 20
		}
		v.Write([]byte(" - General Chat\n"))
		for i := (len(s.messages) - numMessages); i < len(s.messages); i++ {
			m := s.messages[i]
			v.Write([]byte(m.formatMessage()))
			v.Write([]byte("\n"))
		}
	case "right":
		// Refresh the right panel (User List)
		for _, client := range s.clients {
			v.Write([]byte(fmt.Sprintf(" - %s\n", client.name)))
		}
	default:
		return
	}
}

/* var currentChatroom string */

/* func JoinChatroom(gui *gocui.Gui, chatroom string) {
	// Update the current chatroom
	currentChatroom = chatroom

	// Update the left panel to show the current chatroom's information
	gui.Update(func(g *gocui.Gui) error {
		chatView, _ := g.View("left")
		chatView.Clear()
		chatView.Write([]byte(fmt.Sprintf("You joined the %s chatroom!\n", chatroom)))
		// TO DO: display log and room
		return nil
	})
}

func ViewLogs(gui *gocui.Gui, chatroom string) {
	// Update the left panel to show the logs for the selected chatroom
	gui.Update(func(g *gocui.Gui) error {
		chatView, _ := g.View("left")
		chatView.Clear()
		chatView.Write([]byte(fmt.Sprintf("Logs for %s:\n", chatroom)))
		// TO DO: fetch logs for current room

		return nil
	})
} */
