# QuickDB: In-Memory Database

QuickDB is a lightweight, in-memory database written in Go, designed to demonstrate the fundamentals of building a networked database server from scratch. Featuring a custom protocol for efficient client-server communication, QuickDB offers basic functionalities like SET, GET, and DELETE operations on stored key-value pairs. This project serves as a showcase of my journey in mastering Go, particularly focusing on network programming, concurrency, and the implementation of serialization mechanisms.

## Features

- **In-Memory Key-Value Store**: Quick storage and retrieval of data in memory for fast access.
- **Custom Protocol**: A simple and efficient custom binary based network protocol for client-server communication.
- **Concurrent Access**: Utilizes Go's concurrency model to handle multiple client connections simultaneously.
- **Modular Design**: Clean and maintainable code structure, demonstrating best practices in Go development.

## Running QuickDB Server

To start the QuickDB server on the default port (8080):

```bash
go run cmd/quickdb-server/main.go
