package main

import (
	"github.com/jithin-kg/go-in-memory-db/internal/db"
	"github.com/jithin-kg/go-in-memory-db/internal/network"
)

func main() {
	kvStore := db.NewKeyValueStore()
	server := network.NewServer(kvStore)
	server.Start("8080")
}
