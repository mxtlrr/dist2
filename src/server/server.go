package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type ClientData struct {
	clientAddress string // We'll have to connect to the client's HTTP server to send data
	clientPort    string // Connect to this port on the HTTP server.
	threadCount   int    // Amount of threads
	status        int    // Client status
}

const (
	digits int = 20
)

type SortedClientThreads struct {
	clientId    int
	threadCount int
}

var (
	clients          []ClientData
	offset           int64 // Integer offset for computation of digits
	digitsCalculated strings.Builder

	// Priority client IDs
	priorityIDs []SortedClientThreads
)

func main() {
	// Handlers
	http.HandleFunc("/register", registerClient)
	http.HandleFunc("/data", data) // Recieving data from the client
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

func validateClient(c ClientData, ipAddrPort []string) bool {
	return c.clientAddress == ipAddrPort[0] && c.clientPort == ipAddrPort[1]
}

func sortThreadIds() []SortedClientThreads {
	var z []SortedClientThreads

	for i := range clients {
		z = append(z, SortedClientThreads{i, clients[i].threadCount})
	}

	sort.Slice(z, func(i, j int) bool {
		return z[i].threadCount > z[j].threadCount
	})

	return z
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
	for i := range clients {
		if !validateClient(clients[i], ipAddrTmp) {
			io.WriteString(w, "Sorry, you aren't registered. Please register yourself before continuing. Goodbye.")
			return
		}
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
			// Check if priority.
			if len(clients) >= 6 {
				// Only do this if we have enough. If not enough then just
				// prioritize computation

				// Find the index of the current client id
				var curr_index int = 0
				for z := range priorityIDs {
					if priorityIDs[z].clientId == client_id {
						curr_index = z
						break
					}
				}

				// If we're in the lower half, then we should prioritize
				// checking the previous values.
				if curr_index > (len(clients) / 2) {
					// We only want to check three digits, to prevent double
					// computation.

					// Check if computed digits exceeds boundary for amount of
					// digits we want to calculate. There's something we want to
					// actually do!
					if len(digitsCalculated.String()) >= digits {
						io.WriteString(w, fmt.Sprintf("CHECK 3 FROM %d", offset))
					}
				} else {
					// Otherwise, focus on computation
					io.WriteString(w, fmt.Sprintf("COMP %d OFFSET %d", digits, offset))
					offset += 20
				}
			} else {
				// Compute 20 digits of some number
				io.WriteString(w, fmt.Sprintf("COMP %d OFFSET %d", digits, offset))
				offset += 20
			}
			break
		case 1: // Busy
			break
		}
	} else if d_type == "data" {
		d_client := query.Get("data")
		log.Printf("got some data from client %d", client_id)

		digitsCalculated.WriteString(d_client)
		io.WriteString(w, "OK")
	} else if d_type == "check" {
		// This will return OK if no modification needs to be done.
		// otherwise, returns the actual changed digits
		ret := query.Get("return")
		if ret != "OK" {
			// Do something... fix the affected digits?
		}
		// Otherwise nothing needs to be done
		log.Printf("\n")
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
	var z ClientData = ClientData{ipAddrTmp[0], ipAddrTmp[1], threads, 0}
	clients = append(clients, z)

	log.Printf("Finished registeration of client. There are now %d clients in the system\n", len(clients))
	log.Println(z)

	// Make sure client acknowledges.
	n := fmt.Sprintf("OK %d", len(clients))
	io.WriteString(w, n)

	// Update based off registered stuff
	priorityIDs = sortThreadIds()
}
