package main

import (
	"flag"
	"fmt"

	"github.com/BON4/timedQ/internal/server"
)

// @title           TimedQ API
// @version         1.0
// @description     This service will store values provided via API up to certain time. If the value has been accessed, expiration time updates. Key-Value stores in binary file with ttlStore package.

// @host      localhost:8080
// @BasePath  /v1
func main() {
	filePath := flag.String("cfg", "", "path to config.yaml")
	flag.Parse()

	if filePath != nil {
		if *filePath == "" {
			fmt.Println("Please, provide path to config.yaml")
		}

		s, err := server.NewServer(*filePath)
		if err != nil {
			fmt.Printf("INIT ERROR: %s", err.Error())
			return
		}

		if err := s.Run(); err != nil {
			fmt.Printf("RUN ERROR: %s", err.Error())
			return
		}
	}
}
