package turnstyle

import (
	"log"
	"net"
	"os"
)

var (
	logger           = log.New(os.Stdout, "turnstyle: ", log.Lshortfile)
	wasStopRequested = false
)

func Listen() {
	logger.Print("Creating listener for port 8080")
	listener, error := net.Listen("tcp", ":8080")
	if error != nil {
		logger.Fatal("Failed to attach listener to port 8080")
	}
	for {
		if wasStopRequested {
			return
		}
		connection, error := listener.Accept()
		logger.Print("Acquired connection on port 8080")
		if error != nil {
			logger.Fatal("Failed accept connection on port 8080")
		}
		go connect(connection)
	}
}

func Stop() {
	logger.Print("Stopping...")
	wasStopRequested = true
}

func connect(inConnection net.Conn) {
	logger.Print("Connecting to port 8081")
	buffer := make([]byte, 4096)
	outConnection, error := net.Dial("tcp", ":8081")
	if error != nil {
		logger.Fatal("Failed to connect to port 8081")
	}
	logger.Print("Writing to port 8081")
	count, error := inConnection.Read(buffer)
	if count > 0 && error == nil {
		logger.Print(string(buffer[0:count]))
		outConnection.Write(buffer[0:count])
	}
	outConnection.Close()
	logger.Print("Closed connection to 8081")
}
