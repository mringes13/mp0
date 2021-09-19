## MPO
MP0 is a simple program that uses Go routines and Go channels to ping user desired websites.


## How to Run
### Step 1: Initialize analytical ping process
Start the analysis with 'go run mastping.go'
### Step 3: Interact with Command Line
#### Client Process
- To quit while being asked for input, type 'q' and return to terminate the process.
- To quit while the program is running, kill the process by interacting with the system's command line (i.e. OS - Control + C).
- Type the desired websites to ping and analyze, each separated by a space.

## Screenshots


## Workflows


## Custom Data Structures

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


## Exit Codes:

## References


