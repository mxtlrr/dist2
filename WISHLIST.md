# dist2 techology wishlist
Below is a list of stuff that I'd ideally like to add to dist2 to either
- Improve QoL
- Add performance to either
    - Server
    - Client

They are sorted by both what they concern, and the priority of how much
value they add, 10 being "super valuable", 1 being "can be added later on".
***This list is also subject to change.***

## Server
### Fault Tolerance
**Priority**: 8

**Description**: dist2 should not stop running/have clients be desynced
in the case that a client node goes offline for whatever reason.

### Compress Digits to Save on Disk Space
**Priority**: 6

**Description**: Currently, one digit is one byte of disk space. You'd need
around 100 gb for 100 billion digits, which is nowhere near the world
record.

To store more than the current world record amount of digits, which is
[20 trillion](http://www.numberworld.org/y-cruncher/), you'd need 20 TB of
storage, which is viable for large scale servers, but not so viable for
small scale servers, i.e. a home computer, which is what dist2 has been
tested on.

## Client

### Rewrite in Non-Interpreted Language
**Priority**: 4

**Description**: Self-explanatory. Potential speed-ups, but would require a few
different libraries:
- GMP
- Some form of HTTP library (to connect to the server)

Additionally, Pyython has something known as the [Global Interpreter Lock](https://realpython.com/python-gil/).
Essentially, this makes Python code restricted to one thread, but not one process -- which is what dist2 uses
to achieve "multi-threading" (see [issue 4](https://github.com/mxtlrr/dist2/issues/4)). It's "psuedo-multithreading",
it spawns multiple child processes to perform the same task it would with different threads.

Another thing that would benefit the client would be the use of SSE/SIMD instructions. But that comes later, once
the rewrite is done.
