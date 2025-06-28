package server

import (
	"log"
	"os"
	"strconv"
)

func GetPort() int {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		return 3000 // default
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid PORT value: %v", err)
	}
	return port
}
