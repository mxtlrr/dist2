package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
	"strconv"
	"strings"

	"github.com/mxtlrr/dist2/src/server/tdc"
)

type ClientData struct {
	clientAddress string // We'll have to connect to the client's HTTP server to send data
	clientPort    string // Connect to this port on the HTTP server.
	threadCount   int    // Amount of threads
	status        int    // Client status
	currentOffset int64
}

const (
	digits              int   = 20
	MAX_DIGITS_COMPRESS int   = (1<<31)-1 // How many digits before compression?
)

var (
	clients       []ClientData
	offset        int64 = 0 // Integer offset for computation of digits
	totalComputed int64 = 0
	shouldRun     bool  = true
	toCompress    bool  = false // If we compute over some number, then we should
	// compress to save space.
	CSVVals []CSVValue
	
	outFile *os.File
	cT string  // Start time.
	eT string  // End time
	started bool = false
)

func main() {
	log.Println("Welcome to dist2. Parsing csv...")

	// Parse CSV
	csvValTmp, err := os.ReadFile("config.csv")
	if err != nil {
		log.Fatalln(err)
	}

	CSVVals = parseCSV(string(csvValTmp))
	digitsCom, _ := strconv.Atoi(CSVVals[0].value)
	log.Printf("Computing %d digits.", digitsCom)
	if (digitsCom) > MAX_DIGITS_COMPRESS {
		toCompress = true
	}

	outFile, err = os.OpenFile(CSVVals[2].value, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer outFile.Close()

	// Handlers
	http.HandleFunc("/register", registerClient) // Client needds to register itself to prevent unauthorized access
	http.HandleFunc("/data", data)               // Recieving data from the client
	http.HandleFunc("/setstatus", setstatus)

	log.Println("Server running on port 8080")
	go http.ListenAndServe(":8080", nil) // Run in seperate thread so we can do stuff.

	// Just wait.
	for shouldRun {
		continue
	}

	// Save file
	log.Printf("Saving to %s\n", CSVVals[2].value)

	outFile.WriteString(fmt.Sprintf("\n\nComputation started at: %s\nComputation ended at: %s\n\nDist2 v0.0.1\n", cT, eT))
	if toCompress {
		outFile.WriteString("The digits are compressed to save space. Use util/decode to decode the value.\n")
	}
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

		// Write the digits to the file
		if !toCompress {
			if _, err := outFile.WriteString(d_client); err != nil {
				panic(err)
			}
		} else {
			n := tdc.TDCEncodeString(d_client)
			if _, err := outFile.Write(n); err != nil {
				panic(err)
			}
		}
		totalComputed += 20
		io.WriteString(w, "OK")
	}

	// If we've reached the limit, stop executing
	jz, _ := strconv.Atoi(CSVVals[0].value)
	if totalComputed >= int64(jz) {
		eT = time.Now().Format("Jan 2, 2006 15:04:05")
		shouldRun = false
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

	if len(clients) == 1 && !started {
		cT = time.Now().Format("Jan 2, 2006 15:04:05")
		started = true
	}
}




// CSV stuff. To little code to put it in its own thing.
// Maybe i'll do that at some point
type CSVValue struct {
	column string
	value  string
}

func parseCSV(text string) []CSVValue {
	var j []CSVValue
	lines := strings.Split(text, "\n")

	columns := strings.Split(lines[0], ",")
	values := strings.Split(lines[1], ",")

	for value := range values {
		j = append(j, CSVValue{columns[value], values[value]})
	}
	return j
}
