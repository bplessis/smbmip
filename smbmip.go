package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
	"os"
	"strconv"
)

func main() {
	port := flag.Int("port", 8080, "tcp port to listen to")

	flag.Parse()

	db, err := geoip2.Open("/var/lib/GeoIP/GeoLite2-City.mmdb")
	if err != nil {
		fmt.Println("Unable to init GeoDB, GeoTrace disabled:", err.Error())
		db = nil
	} else {
		defer db.Close()
	}

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
		go handleConnection(conn, db)
	}
}

type GeoLoc struct {
	Country string `json:"country"`
	City    string `json:"city"`
	State   string `json:"state"`
}
type Connection struct {
	Name string `json:"marbuk"`
	IP   net.IP `json:"ip"`
	Geo  GeoLoc `json:"geoloc"`
}

// Handles incoming requests.
func handleConnection(conn net.Conn, db *geoip2.Reader) {
	var thisConn Connection

	thisConn.Name = "sbmip"

	thisConn.IP = conn.RemoteAddr().(*net.TCPAddr).IP

	if db != nil {
		// If you are using strings that may be invalid, check that ip is not nil
		record, err := db.City(thisConn.IP)
		if err != nil {
			fmt.Println("GeoLoc Failed: " + err.Error())
		} else {
			thisConn.Geo.City = record.City.Names["en"]
			thisConn.Geo.Country = record.Country.Names["en"]
			//thisConn.Geo.State = record.Subdivisions[0].Names["en"]
		}
	}

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
