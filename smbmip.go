package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	port := flag.Int("port", 8080, "tcp port to listen to")

	flag.Parse ()
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		// handle error
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Listening on port " + strconv.Itoa(*port))
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("Error accepting connection : " + err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

type Connection struct {
	Name        string      `json:"marbuk"`
	IP          net.IP      `json:"ip"`
}

// Handles incoming requests.
func handleConnection(conn net.Conn) {
	var thisConn Connection

	thisConn.Name = "sbmip"

	thisConn.IP = conn.RemoteAddr().(*net.TCPAddr).IP

	// Send a response back to person contacting us.
	b, err := json.Marshal(thisConn)
	if err != nil {
		conn.Write([]byte("Message received from " + thisConn.IP.String()))
		fmt.Println("JSON Conversion error: " + err.Error())
	} else {
		conn.Write(b)
	}

	// Close the connection when you're done with it.
	conn.Close()
}
