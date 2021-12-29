# tcp-replay
Replay tcp content to the destination address.

## Usage
Read data from stdin, which may be filled by a pipe between tcpdump.

## Example
1. `cat data.pcap | tcp-replay -t localhost:8080 -d 1s`
2. `tcpdump -i eth0 -U -w - tcp and dst port 80 | tcp-replay -t localhost:8080 -d 1s`
