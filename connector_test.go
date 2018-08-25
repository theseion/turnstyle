package turnstyle_test

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"turnstyle"
)

var (
	logger   = log.New(os.Stdout, "turnstyle-test: ", log.Lshortfile)
	listener net.Listener
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	setUp()
	defer tearDown()
	os.Exit(m.Run())
}

func setUp() {
	go func() {
		listener, error := net.Listen("tcp", ":8081")
		if error != nil {
			panic("Could not attache listener to port 8081")
		}
		inConnection, error := listener.Accept()
		if error != nil {
			logger.Fatal("Failed accept connection on port 8080")
		}
		defer inConnection.Close()
		request, error := bufio.NewReader(inConnection).ReadString('\n')
		logger.Print(request)
		if error != nil && request != "GET / HTTP/1.0\r\n\r\n" {
			panic("Could not read from connection")
		}
	}()
}

func tearDown() {
	listener.Close()
	turnstyle.Stop()
}

func TestListening(t *testing.T) {
	go turnstyle.Listen()
	logger.Print("dialling")
	var outConnection net.Conn
	var err error
	for outConnection, err = net.Dial("tcp", ":8080"); err != nil; {
	}
	defer outConnection.Close()
	logger.Print("writing")
	fmt.Fprintf(outConnection, "GET / HTTP/1.0\r\n\r\n")
}

// func TestSimpleProxy(t *testing.T) {
// 	go turnstyle.Listen()
// 	go func() {
// 		listener, error := net.Listen("tcp", ":8081")
// 		if error != nil {
// 			t.FailNow()
// 		}
// 		inConnection, error := listener.Accept()
// 		if error != nil {
// 			logger.Fatal("Failed accept connection on port 8080")
// 		}
// 		request, error := bufio.NewReader(inConnection).ReadString('\n')
// 		if error != nil && request != "GET / HTTP/1.0\r\n\r\n" {
// 			t.FailNow()
// 		}
// 	}()
// 	time.Sleep(100 * time.Millisecond)
// 	logger.Print("dialling")
// 	outConnection, error := net.Dial("tcp", ":8080")
// 	if error != nil {
// 		t.FailNow()
// 	}
// 	logger.Print("writing")
// 	fmt.Fprintf(outConnection, "GET / HTTP/1.0\r\n\r\n")
// 	logger.Print("reading")
// 	response, error := bufio.NewReader(outConnection).ReadString('\n')
// 	logger.Print(response)
// 	if error != nil && response != "success\r\n\r\n" {
// 		t.FailNow()
// 	}
// }
