package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ClientData struct {
	clientAddress string // We'll have to connect to the client's HTTP server to send data
	clientPort    string // Connect to this port on the HTTP server.
	threadCount   int    // Amount of threads
	status        int    // Client status
	currentOffset int64
}

const (
	digits int = 20
)

var (
	clients          []ClientData
	offset           int64 // Integer offset for computation of digits
	digitsCalculated strings.Builder
)

func main() {
	// Handlers
	http.HandleFunc("/register", registerClient) // Client needds to register itself to prevent unauthorized access
	http.HandleFunc("/data", data)               // Recieving data from the client
	http.HandleFunc("/setstatus", setstatus)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setstatus(w http.ResponseWriter, r *http.Request) {
	query, _ := url.ParseQuery(r.URL.RawQuery)
	status, _ := strconv.Atoi(query.Get("val"))
	client_id, _ := strconv.Atoi(query.Get("client_id"))

	clients[client_id].status = status
}

func data(w http.ResponseWriter, r *http.Request) {
	// Check if the IP address actually corresponds to a client
	ipAddrTmp := strings.Split(r.RemoteAddr, ":")
	var good bool = false
	for z := range clients {
		c := clients[z]
		if ipAddrTmp[0] == c.clientAddress {
			if ipAddrTmp[1] == c.clientPort {
				good = true
			}
		}
	}

	if !good {
		io.WriteString(w, "Sorry, validation failed!")
		return
	}

	log.Println("Incoming data / request from client!")
	query, _ := url.ParseQuery(r.URL.RawQuery)
	d_type := query.Get("type")
	client_id, _ := strconv.Atoi(query.Get("client_id"))

	fmt.Println(query)
	log.Printf("Got \"%s\" as the type.", d_type)

	// TODO: do something with type given to server.

	if d_type == "request" {
		switch clients[client_id].status {
		// Ready
		case 0:
			// Offsetted?
			for n := range clients {
				if n+1 != len(clients) {
					// Fix bug where clients would be desynchronized and compute
					// over and over again
					if clients[n+1].currentOffset < clients[n].currentOffset {
						fmt.Println("Uh-oh! bad thing happen.")
						clients[n+1].currentOffset += int64(digits) * int64(len(clients))
					}
				}
			}
			// Compute 20 digits of some number
			io.WriteString(w, fmt.Sprintf("COMP %d OFFSET %d", digits, offset))
			offset += 20
			clients[client_id].currentOffset = offset
		}
	} else if d_type == "data" {
		d_client := query.Get("data")
		log.Printf("got some data from client %d", client_id)

		digitsCalculated.WriteString(d_client)
		io.WriteString(w, "OK")
	}
}

func registerClient(w http.ResponseWriter, r *http.Request) {
	ipAddrTmp := strings.Split(r.RemoteAddr, ":")
	log.Println("Recieving client registration")

	// Get threads
	query, _ := url.ParseQuery(r.URL.RawQuery)
	threads, e := strconv.Atoi(query.Get("threads"))

	if e != nil {
		log.Printf("Error registering! %s is not base-10. Cannot register client.", query.Get("threads"))
		io.WriteString(w, "FAILURE")
		return
	}

	// Server stuff
	log.Printf("Threads: %d | IP Address: %s\n", threads, ipAddrTmp[0])

	// Add client data to the clients
	var z ClientData = ClientData{ipAddrTmp[0], ipAddrTmp[1], threads, 0, 0}
	clients = append(clients, z)

	log.Printf("Finished registeration of client. There are now %d clients in the system\n", len(clients))
	log.Println(z)
	log.Println(clients)

	// Make sure client acknowledges.
	n := fmt.Sprintf("OK %d", len(clients))
	io.WriteString(w, n)
}
