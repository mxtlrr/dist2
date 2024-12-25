# dist2
Parallelized computation of digits from $\sqrt{2}$

## What does dist2 collect?
Due to dist2 being parallelized and having to be connected to a server
to get/perform calculations, all that dist2 collects is:
- Your IP address (to send/recieve information)
- Processor information (amount of threads)

That's all, nothing more, nothing less. Unfortunately, there is no way to
get around this. dist2 needs to know your CPU/threads amount to allocate
more work to clients with more computational strength.


# Compilation
The server and client are written in different programming languages,
meaning you'll need to get multiple dependencies:
- Golang (server)
- Python (client)

## Sever
As stated above, the server is written in Golang. You can compile it with
```
go build server/server.go
```

## Client
The client is written in Python. Just run `python client.py`.

# Usage

## Client
The client uses an INI file for locating where the server is. The INI file
looks like this:
```ini
[config]
svip = '127.0.0.1'
port = 8080
```

Modify these as you wish and then run the client. It will fail if the server
is not running.

## Server
As of right now, the client runs on port 8080 on localhost, which is what the client connects to. You'll
have to do slight configuration: dist2 uses a [csv file](./src/server/config.csv) for configuration (it's
easier than ini)

This needs to be in the same folder as the binary of the server. Currently, there are only two values you need
to modify: `limit_digits` and `output_file`. It is self explanatory. Once you launch, if all goes well, you should see
something similar in your terminal:
```
2024/12/14 10:18:20 Welcome to dist2. Parsing csv...
2024/12/14 10:18:20 Computing 500 digits.
2024/12/14 10:18:20 Server running on port 8080
```

# Records
| Digit Count | Time Taken | Number of Clients | Date |
| ----------- | ---------- | ----------------- | ---- |
| 10,000      | 7  seconds | 2   | 15 December 2024   |

To submit records, use the issues tab with the 'record' tag, and attach the `digits.txt` file
generated.

# Architecture
1. Upon starting a client, the client sends information about itself.
2. The server will acknowledge this, and waits for the client to send a RDY
message.
3. Upon the RDY message, the server sends a data packet to the client, to compute $n$ digits at an offset $o$, which starts at 0 digits past the decimal points.
	- On other clients, it will tell the client to compute $k+n$ digits, aka, computing digits sequentially.
4. Once computation is completed, the client sends a DONE message. 
5. This process repeats until all digits have been computed. Once that's
done, all clients go into validaton mode.
6. Validation mode essentially consists of going through the file again and
checking each section of digits against a different algorithm, and replacing
the invalid digits in the file if needed.
7. Step 6 goes through over all digits in the file. Once that's done the
server exits.

You can view the diagram in the [documentation](./docs/graph.md) directory.
Additionally, in the same directory is an [explanation of how the algorithms work](./docs/algs.md).

## Architecture Implementation Roadmap
- [X] Send information to server about client information
- [ ] Server acknowledges and waits for client to be ready to send info
	- [X] Tells the client to compute $n$ digits at offset $o$.
        - [ ] If a client has more threads, then make it compute more digits than other
        clients.
	- [ ] Prioritizes clients with more than $j$ threads
	- [X] Client can both
		- [X] Compute $\sqrt{2}$ with decent precision, enough for what
		- [X] Send the data back to the server, which stores it.
- [X] Implement checking
	- [X] Server-side signal
	- [X] Client-side
- [X] Optional stuff
	- [X] Send a terminate signal to all clients once a desired precision
is wanted

## General Roadmap
- [ ] Implement the entire architecture
- [ ] Implement things from [the technological wishlist](./WISHLIST.md)

## Notes
1. "Weaker" clients mean that they don't have as much computing power. If client $A$ has 4 threads, and client $B$ has
2 threads, then client $A$ will have priority over computation, and client $B$ will have priority over compiling.
