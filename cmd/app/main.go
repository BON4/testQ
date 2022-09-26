package main

import (
	"flag"
	"fmt"

	"github.com/BON4/timedQ/internal/server"
)

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
