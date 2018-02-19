package main

import (
	"strconv"
	"fmt"
	"net"
	"os"
)

func main() {
	var port int = 8080 ;

	ln, err := net . Listen("tcp", ":" + strconv.Itoa( port ) )
	if err != nil {
		// handle error
		fmt . Println("Error listening:", err.Error())
		os . Exit(1)
	}
	defer ln . Close ()

	fmt . Println ( "Listening on port " + strconv.Itoa ( port ) )
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt . Println ( "Error accepting connection : " + err.Error () )
			os . Exit (1)
		}
		go handleConnection(conn)
	}
}

// Handles incoming requests.
func handleConnection(conn net.Conn) {
	// Send a response back to person contacting us.
	var ls = conn . LocalAddr () . (*net.TCPAddr)
	var lc = conn . RemoteAddr () . (*net.TCPAddr)

	conn . Write( []byte("Message received from " + lc.IP.String() + " to " + ls.IP.String() ) )
	// Close the connection when you're done with it.
	conn . Close()
}
