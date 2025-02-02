# Description

This project consists on recreating the NetCat in a Server-Client Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

The following features are implemented:

- TCP connection between server and multiple clients (relation of 1 to many).
- A name requirement to the client.
- Control connections quantity.
- Clients must be able to send messages to the chat.
- Do not broadcast EMPTY messages from a client.
- Messages sent, are identified by the time that was sent and the user name of who sent the message, example :    
`[2020-01-20 15:48:41][client.name]:[client.message]`
- If a Client joins the chat, all the previous messages sent to the chat must be uploaded to the new Client.
- If a Client connects to the server, the rest of the Clients must be informed by the server that the Client joined the group.
- If a Client exits the chat, the rest of the Clients must be informed by the server that the Client left.
- All Clients must receive the messages sent by other Clients.
- If a Client leaves the chat, the rest of the Clients must not disconnect.
- If there is no port specified, then set as default the port 8989. Otherwise, program will respond with usage message:   
`[USAGE]: ./TCPChat $port`
- Gui provided on the server terminal that shows active users, messages on chat and allows the admin to send messages to the chat
- Log file of each chat is saved on hard drive as .log file named with the date and time of creation

# Usage

To run the server: `go run . $port`
To connect as a client: `nc localhost $port` no `:` prefix required. The default port is `:8989`

# Dependancies

- go lang 1.23.4
- github.com/jroimartin/gocui v0.5.0
