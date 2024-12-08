package main

import (
	"io"
	"log"
	"strconv"
	"strings"
	"net/http"
	"net/url"
)

type ClientData struct {
	clientAddress		string			// We'll have to connect to the client's HTTP server to send data
	clientPort			string			// Connect to this port on the HTTP server.
	threadCount			int					// Amount of threads
}

var (
	clients []ClientData
)

func main() {
	http.HandleFunc("/register", registerClient)
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


// Remove the port
func TruncIPAddr(fullStr string) string {
	return strings.Split(fullStr, ":")[0]
}

func registerClient(w http.ResponseWriter, r *http.Request){
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
	log.Printf("Threads: %d | IP Address: %s\n", threads, TruncIPAddr(r.RemoteAddr));

	// Add client data to the clients
	var z ClientData = ClientData{TruncIPAddr(r.RemoteAddr), strings.Split(r.RemoteAddr, ":")[1], threads}
	clients = append(clients, z)

	log.Printf("Finished registeration of client. There are now %d clients in the system\n", len(clients))
	log.Println(z)
}
