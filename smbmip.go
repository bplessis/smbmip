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

type GeoLoc struct {
	Country string `json:"country"`
	City    string `json:"city"`
	State   string `json:"state"`
	ASN     string `json:"AS"`
	ASOrg   string `json:"AS_Org_Name"`
}
type Connection struct {
	Name string `json:"marbuk"`
	IP   net.IP `json:"ip"`
	Geo  GeoLoc `json:"geoloc"`
}

var dbCity *geoip2.Reader
var dbASN *geoip2.Reader

func main() {
	var err error

	port := flag.Int("port", 8080, "tcp port to listen to")
	geoliteCityDB := flag.String("geocity", "/var/lib/GeoIP/GeoLite2-City.mmdb", "Path to GeoLite2 City Database")
	geoliteASNDB := flag.String("geoasn", "/var/lib/GeoIP/GeoLite2-ASN.mmdb", "Path to GeoLite2 ASN Database")

	flag.Parse()

	dbCity, err = geoip2.Open(*geoliteCityDB)
	if err != nil {
		fmt.Println("Unable to init GeoDB, GeoTrace disabled:", err.Error())
		dbCity = nil
	} else {
		defer dbCity.Close()
	}

	dbASN, err = geoip2.Open(*geoliteASNDB)
	if err != nil {
		fmt.Println("Unable to init GeoDB, GeoTrace disabled:", err.Error())
		dbASN = nil
	} else {
		defer dbASN.Close()
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
		go handleConnection(conn)
	}
}

// Handles incoming requests.
func handleConnection(conn net.Conn) {
	var thisConn Connection

	thisConn.Name = "sbmip"

	thisConn.IP = conn.RemoteAddr().(*net.TCPAddr).IP

	if dbCity != nil {
		// If you are using strings that may be invalid, check that ip is not nil
		record, err := dbCity.City(thisConn.IP)
		if err != nil {
			fmt.Println("GeoLoc Failed: " + err.Error())
		} else {
			thisConn.Geo.City = record.City.Names["en"]
			thisConn.Geo.Country = record.Country.Names["en"]
			if len(record.Subdivisions) > 0 {
				thisConn.Geo.State = record.Subdivisions[0].Names["en"]
			}
		}
	}
	if dbASN != nil {
		// If you are using strings that may be invalid, check that ip is not nil
		record, err := dbASN.ASN(thisConn.IP)
		if err != nil {
			fmt.Println("ASN GeoLoc Failed: " + err.Error())
		} else {
			thisConn.Geo.ASN = "AS" + strconv.FormatUint(uint64(record.AutonomousSystemNumber), 10)
			thisConn.Geo.ASOrg = record.AutonomousSystemOrganization
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
