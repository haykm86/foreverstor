# ForeverStor

ForeverStor is a peer-to-peer file storage prototype written in Go. Each node runs a `FileServer` that can join a TCP-based overlay network, replicate encrypted file blobs, and respond to fetch requests from peers. The project bundles a simple content-addressable storage layer, symmetric encryption helpers, and a lightweight TCP transport abstraction so the system can be extended or embedded in other applications.

## Highlights

- Peer discovery via static bootstrap list and outbound connection attempts.
- Streaming transport built on TCP with a tiny RPC envelope defined in `p2p`.
- Pluggable storage with configurable path transformers for content-addressable layouts.
- AES-CTR encryption of network payloads so replicated data is never stored in cleartext.
- Hash-derived file identifiers and per-node storage roots for collision avoidance.

## How It Works

1. **Transport layer (`p2p/`):** `TCPTransport` wraps Go's `net` package, handles handshakes, and multiplexes messages/streams through an internal RPC channel. A minimal `DefaultDecoder` inspects the first byte to differentiate control messages from raw streams.
2. **Storage (`storage.go`):** Files are split across a deterministic directory tree derived from a SHA-1 hash. The `Store` type handles reads, writes, deletions, and optional decryption before persisting the payload.
3. **Server logic (`server.go`):** `FileServer` tracks peers, broadcasts store/get requests, and coordinates encryption, replication, and local caching. Nodes negotiate file transfers by announcing intent (`MessageStoreFile`, `MessageGetFile`) before streaming bytes.
4. **Crypto helpers (`crypto.go`):** Utility functions provide node IDs, hash keys with MD5 (suitable for demo purposes), and encrypt/decrypt streams with AES-CTR.

The demo program in `main.go` spins up three nodes, replicates an example object, deletes one replica, and finally reads the blob back from the network.

## Getting Started

```bash
go run .
```

By default the sample bootstraps three nodes on `:33033`, `:44044`, and `:55055`, each writing data under `<listenAddr>_network/`. Adjust the ports or bootstrap list inside `main.go` to experiment with alternative topologies.

To experiment with different network topologies today, edit the `makeServer` calls in `main.go` and relaunch the program. Wiring a CLI parser around those options is a natural next step if you want to expose runtime configuration.

## Project Layout

- `main.go` – demo entrypoint that wires three servers together.
- `server.go` – core file server implementation, message handling, and peer bookkeeping.
- `storage.go` – content-addressable storage primitives and helpers.
- `crypto.go` – symmetric encryption utilities plus ID/hash helpers.
- `p2p/` – TCP transport, peer abstraction, and decoding logic.
- `*_test.go` – store, crypto, and transport unit tests.

## Testing and Tooling

Run the full test suite with:

```bash
go test ./...
```

`go.mod` targets Go 1.22 and depends on `github.com/stretchr/testify` for assertions.

## Ideas for Extension

- Add a CLI layer to configure listen addresses, bootstrap peers, and local storage roots.
- Replace the in-memory bootstrap list with a gossip layer or DHT lookup.
- Harden cryptography (e.g., authenticated encryption, rotating keys, envelope encryption per object).
- Add integrity checks around stored objects (e.g., compare SHA-256 digest after replication).
- Persist peer metadata and retry logic to survive transient network failures.

ForeverStor is intentionally compact, making it a good foundation for exploring distributed storage concepts or experimenting with transport abstractions in Go.
