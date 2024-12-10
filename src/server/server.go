package main

import (
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
}

var (
	clients []ClientData
)

func main() {
	// Handlers
	http.HandleFunc("/register", registerClient)
	http.HandleFunc("/data", data) // Recieving data from the client

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
The client will send a variety of messages to the server,
* telling it if/when it is ready, and data from various tasks.
* You can read the architecture about it. Once the client
* configures itself and registers itself on the server, it sends
* a RDY to the server, i.e.

*  /data?type="RDY", then the server will tell it what to do
*/
func data(w http.ResponseWriter, r *http.Request) {
	// Check if the IP address actually corresponds to a client
	ipAddrTmp := strings.Split(r.RemoteAddr, ":")
	var (
		ipAddr      string = ipAddrTmp[0]
		port        string = ipAddrTmp[1]
		validClient bool   = false
	)

	for i := range clients {
		if clients[i].clientAddress == ipAddr && clients[i].clientPort == port {
			validClient = true
			break
		}
	}

	if validClient == false {
		io.WriteString(w, "Sorry, you aren't registered. Please register yourself before continuing. Goodbye.")
		return
	}

	log.Println("Incoming data / request from client!")
	query, _ := url.ParseQuery(r.URL.RawQuery)
	d_type := query.Get("type")

	log.Printf("Got \"%s\" as the type.", d_type)
	// TODO: do something with type given to server.
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
	var z ClientData = ClientData{ipAddrTmp[0], ipAddrTmp[1], threads}
	clients = append(clients, z)

	log.Printf("Finished registeration of client. There are now %d clients in the system\n", len(clients))
	log.Println(z)

	// Make sure client acknowledges.
	io.WriteString(w, "OK")
}
