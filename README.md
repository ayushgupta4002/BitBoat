# Distributed Cache System

## Overview

This project implements a distributed cache system using Go. It is designed to handle caching efficiently across multiple nodes, providing capabilities for storing, retrieving, and managing key-value pairs. The system includes features like TTL (time-to-live) for cache entries, distributed node communication, and administrative functionality for managing cluster operations.
![image](https://github.com/user-attachments/assets/9264e2b7-8036-4a3b-b24d-466efe6b6555)


---

## Features

- **Distributed Architecture**: The cache system can be deployed across multiple nodes.
- **TTL Support**: Cached entries expire after a configurable time-to-live duration.
- **Cluster Management**: Administrative functionality for managing nodes and subscribers.
- **Key Operations**: Supports `SET`, `GET`, `DELETE`, and `HAS` operations.
- **Scalable Design**: Nodes can join or leave the cluster dynamically.
- **Client Support**: A client library is provided for interacting with the cache system programmatically.

---

## Getting Started

### Prerequisites

- Go 1.18+ installed on your system.

### Installation

Clone the repository:

```bash
git clone https://github.com/ayushgupta4002/bitboat.git
cd bitboat
```

### Running the Server

Start a Admin server node:

```bash
go run main.go --listenaddr localhost:8080
```

Start a Subscriber node that connects to Admin:

```bash
go run main.go --listenaddr localhost:8081 --adminaddr localhost:8080
```

![image](https://github.com/user-attachments/assets/f82ef0f7-5993-4354-9bfe-1116669c00cf)

---

## Client Usage

The client library allows applications to interact with the cache system. Clients can interact with any Node ( Subscriber or Admin ), Example usage:

```go
package main

import (
	"context"
	"log"

	"github.com/ayushgupta4002/bitboat/client"
)

func main() {
	client, err := client.NewClient("localhost:8080", client.ClientOpts{}) // [ here specify address of subscriber or admin node as per client need]
	if err != nil {
		log.Fatal("Error connecting to cache:", err)
	}

	// Set a key
	err = client.Set(context.Background(), []byte("key1"), []byte("value1"), 3000000000)
	if err != nil {
		log.Fatal("Error setting value:", err)
	}

	// Get a key
	value, err := client.Get(context.Background(), []byte("key1"))
	if err != nil {
		log.Fatal("Error getting value:", err)
	}
	log.Println("Value:", string(value))

	// Delete a key
	err = client.Delete(context.Background(), []byte("key1"))
	if err != nil {
		log.Fatal("Error deleting key:", err)
	}

	client.Close()
}
```

---

## Protocol Details

The cache system uses a custom binary protocol for communication between nodes and clients. Commands supported:

1. **SET**: Store a key-value pair with a TTL.
2. **GET**: Retrieve the value of a key.
3. **DELETE**: Remove a key from the cache.
4. **HAS**: Check if a key exists in the cache.
5. **JOIN**: Add a Subscriber node to the Admin cluster.

---

## Project Structure

- **`cache`**: Implements the in-memory cache.
- **`client`**: Provides a Go client library for interacting with the cache system.
- **`proto`**: Defines the custom protocol for communication.
- **`main.go`**: The entry point for starting server and client nodes.

---

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
