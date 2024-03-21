
# Bittorrent Implementation README

## Introduction

Welcome to the Bittorrent Implementation project! This project aims to guide you through the process of building a simplified Bittorrent client in stages, covering essential functionalities step by step.

Ensure you have `go (1.19)` installed locally
Run `./bittorrent.sh` to run your program, which is implemented in
   `cmd/bittorrent/main.go`.




## Phase 1: Bencode Decoding

Bencode is a serialization format used in the BitTorrent protocol. This phase focuses on decoding strings. Here's how to run the program:

```bash
$ ./bittorrent.sh decode 5:hello
```

## Phase 2: Bencoded Integers

Extend the decode command to support bencoded integers:

```bash
$ ./bittorrent.sh decode i52e
```

## Phase 3: Bencoded Lists

Extend the decode command to support bencoded lists:

```bash
$ ./bittorrent.sh decode l5:helloi52ee
```

## Phase 4: Bencoded Dictionaries

Extend the decode command to support bencoded dictionaries:

```bash
$ ./bittorrent.sh decode d3:foo3:bar5:helloi52ee
```

## Phase 5: Torrent File Parsing

Parse a torrent file and print information about the torrent:

```bash
$ ./bittorrent.sh info sample.torrent
```

## Phase 6: Info Hash Calculation

Calculate the info hash for a torrent file:

```bash
$ ./bittorrent.sh info sample.torrent
```

## Phase 7: Piece Length and Hashes

Print piece length and piece hashes in hexadecimal format:

```bash
$ ./bittorrent.sh info sample.torrent
```

## Phase 8: Tracker Interaction

Make a GET request to a HTTP tracker to discover peers:

```bash
$ ./bittorrent.sh peers sample.torrent
```

## Phase 9: Handshake with a Peer

Establish a TCP connection with a peer and complete a handshake:

```bash
$ ./bittorrent.sh handshake sample.torrent <peer_ip>:<peer_port>
```

## Phase 10: Piece Download

Download one piece and save it to disk:

```bash
$ ./bittorrent.sh download_piece -o /tmp/test-piece-0 sample.torrent 0
```

## Phase 11: Full File Download

Download the entire file and save it to disk:

```bash
$ ./bittorrent.sh download -o tmp/test.txt sample.torrent
``` 

Feel free to explore each phase of the project and expand your understanding of the Bittorrent protocol implementation!
