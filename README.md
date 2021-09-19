## MP0
MP0 is a simple application that pings desired websites, as defined by the user, using Go-Routines and Go-Channels. 
This project was assigned to practice the newly learned Go-Routines and Go-Channels.

## How to Run
### Step 1: Initialize ping of websites
Start the analytical ping with `go run mastping.go`
### Step 2: Interact with Command Line

#### User Ping Websites Process
- To quit the program when being asked to input websites, return 'q'.
- To quit the program while running, interact with the system command line (i.e. OS Terminal - Control + C)
- Enter the desired websites to be pinged.

## Screenshots

## Workflows
### Ping Process
![MP0 Workflow Diagram](https://user-images.githubusercontent.com/60116121/133932682-9a37ebe8-20af-487f-95b2-b4035317fc1b.png)


## Custom Data Structures
1. Message Struct
```go
type Message struct {
   To string
   From string
   Content string
}
```
2. Username-Connection Lookup Map
```go
var clientLookup = make(map[string]net.Conn)
```
## Package Design
### Application
- `chatroomparsing.go` contains functions for parsing initial command line arguments and reading the command line to
  terminate a chatroom process upon user request
- `clientparsing.go` contains functions for parsing initial command line arguments and reading the command line to 
  construct messages and terminate a client process upon user request
### Network
- `chatroom.go` contains functions for listening to a TCP port as well as routing messages from client to client
- `client.go` contains functions for establishing a connection to a TCP port as well as sending and receiving messages 
  to and from the chatroom
- `communication.go` contains functions for writing and reading messages to a TCP channel via gob
### Messages
`messages.go` contains the Message struct
### Error Checker
`errorchecker.go` contains a function to check for errors for initial TCP connection functions.
### Images
Contains all the images used in this README.


## Exit Codes:
- `0`: Successful
- `1`: Incorrect command line input format
- `2`: External package function error

## References
- My error checking function, `CheckError()`, is a modified version of sample code from [Network Programming with Go](https://ipfs.io/ipfs/QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/socket/tcp_sockets.html).
- My exact implementation of establishing a TCP connection on both client and server side was taken from [this linode tutorial](https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/).
