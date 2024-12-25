package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

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
	MAX_DIGITS_COMPRESS int  = (1 << 31) - 1 // How many digits before compression?
	MAX_SIZE            int  = 1000          // Max size of the "validSave" string.
	MIN_THREADS         int8 = 2             // Anything leq this will focus on checking.
)

var (
	clients       []ClientData
	toCompress    bool    = false // If we compute over some number, then we should
	shouldRun     bool    = true
	shouldCompute bool    = true
	average       float64 = 0.0
	digits        int     = 100
	clientCount   int     = 0
	offset        int64   = 0 // Integer offset for computation of digits
	totalComputed int64   = 0

	// compress to save space.
	CSVVals []CSVValue

	// Timing and other things
	outFile *os.File
	cT      string // Start time.
	cTime   time.Time
	eTime   time.Time
	chTime  time.Time
	eT      string // End time
	started bool   = false

	// Checking stuff
	digits_for_check   int = 100 // check 100 digits at a time
	offset_digit_check int = 0
	total              int = 0 // Total things changed
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

	outFile.Close()

	// This is the dumbest shit ever. I don't know but it doesn't work if I use outFile. Good enough for now.
	newFileFuckYou, _ := os.OpenFile(CSVVals[2].value, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	newFileFuckYou.WriteString(fmt.Sprintf("\n\nComputation started at: %s\nComputation ended   at: %s\n", cT, eT))
	newFileFuckYou.WriteString(fmt.Sprintf("Validation of %s digits time:       %s\n", CSVVals[0].value, time.Time{}.Add(chTime.Sub(eTime)).Format("15:04:05.000")))
	newFileFuckYou.WriteString(fmt.Sprintf("Computation duration:                %s\n\nDist2 v0.0.1\n", time.Time{}.Add(eTime.Sub(cTime)).Format("15:04:05.000")))
	newFileFuckYou.WriteString(fmt.Sprintf("Total invalid digits: %d\n", total))
	newFileFuckYou.WriteString(fmt.Sprintf("Percentage (wrong):   %.3f%%\n", (float32(total)/float32(digitsCom))*100))
	if toCompress {
		newFileFuckYou.WriteString("The digits are compressed to save space. Use util/decode to decode the value.\n")
	}
	newFileFuckYou.Close()
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

	jz, _ := strconv.Atoi(CSVVals[0].value)
	if d_type == "request" {
		switch clients[client_id].status {
		// Ready
		case 0:
			if shouldCompute {
				// If it took us less than 5 ms to compute, do nothing. Otherwise,
				// double digit count. We'll end up writing a multithreaded thing
				// for this in the client.
				if average <= 0.005 {
					digits *= 2
				}
				// Offsetted?
				for n := range clients {
					if n+1 != len(clients) {
						// Fix bug where clients would be desynchronized and compute
						// over and over again
						if clients[n+1].currentOffset < clients[n].currentOffset {
							clients[n+1].currentOffset += int64(digits) * int64(len(clients))
						}
					}
				}

				io.WriteString(w, fmt.Sprintf("COMP %d OFFSET %d MAX %d", digits, offset, jz))
				offset += int64(digits)
				clients[client_id].currentOffset = offset
			} else { // All clients start now checking
				// Open the file
				file, err := os.Open(CSVVals[2].value)
				if err != nil {
					panic(err)
				}

				// Seek to the desired offset
				_, err = file.Seek(int64(offset_digit_check), io.SeekStart)
				if err != nil {
					panic(err)
				}
				// Read `n` bytes
				buf := make([]byte, digits_for_check)
				if _, e := file.Read(buf); e != nil {
					panic(e)
				}
				value := string(buf[:])
				fmt.Printf("FUCK YOU!!!!: %d\n", offset_digit_check)
				io.WriteString(w, fmt.Sprintf("CHECK %d OFFSET %d STUFF %s", digits_for_check, offset_digit_check, value))
				offset_digit_check += digits_for_check
				fmt.Printf("FUCK YOU!!!!: %d\n", offset_digit_check)
				file.Close() // prevent any memory leaks
			}
		}
	} else if d_type == "data" {
		vaC := query.Get("typeOfData")
		if vaC == "comp" {
			d_client := query.Get("data")
			kl, _ := strconv.ParseFloat(query.Get("timing"), 64)
			if clientCount < len(clients) {
				average += kl
				clientCount += 1
			} else {
				average = kl / float64((len(clients)))
				clientCount = 0
			}

			log.Printf("got some data from client %d", client_id)
			log.Printf("Average timing for clients is %.6f", average)

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
			totalComputed += int64(digits)
			io.WriteString(w, "OK")
		} else {
			var (
				digits   string = query.Get("digs")
				retVal   string = query.Get("ret_val")
				original string = query.Get("originalData")
			)
			if retVal == "BAD" {
				e := replaceInFile(CSVVals[2].value, original, digits)
				if e != nil {
					log.Fatal(e)
				}
				total += digits_for_check
			}
		}
	}

	// If we've reached the limit, stop computing
	// and start
	if totalComputed >= int64(jz) {
		eTime = time.Now()
		eT = eTime.Format("Jan 02, 2006 15:04:05.000")
		shouldCompute = false
	}

	if !shouldCompute && offset_digit_check >= jz {
		chTime = time.Now()
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
		cTime = time.Now()
		cT = cTime.Format("Jan 2, 2006 15:04:05.000")
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

func replaceInFile(fileName, searchPattern, replacePattern string) error {
	tempFileName := fileName + ".tmp"

	sourceFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open source file '%s': %w", fileName, err)
	}
	defer sourceFile.Close()

	tempFile, err := os.Create(tempFileName)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	reader := bufio.NewReader(sourceFile)
	writer := bufio.NewWriter(tempFile)

	pattern := []byte(searchPattern)
	replacement := []byte(replacePattern)

	for {
		chunk, err := reader.ReadBytes('\n') // Process line-by-line
		if err != nil && len(chunk) == 0 {   // End of file or error
			break
		}

		modifiedChunk := bytes.ReplaceAll(chunk, pattern, replacement)

		if _, writeErr := writer.Write(modifiedChunk); writeErr != nil {
			return fmt.Errorf("failed to write to temporary file: %w", writeErr)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	if err := os.Rename(tempFileName, fileName); err != nil {
		return fmt.Errorf("failed to replace the original file: %w", err)
	}

	return nil
}
