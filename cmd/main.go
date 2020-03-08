package main

import (
	"os"

	"reversi/client"
)

func main() {
	os.Exit(client.NewReversi().Run())
}
