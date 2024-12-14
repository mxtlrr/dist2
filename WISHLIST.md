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
**Description**: Self-explanatory. Potential speed-ups, but would require
an external large number library, i.e. [GMP](https://gmplib.org/)
