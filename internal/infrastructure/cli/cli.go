package cliCleint

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"
)

type CliClient struct {
}

func NewCliClient() *CliClient {

	return &CliClient{}
}

// Run starts reading user input from the command line continuously and sends it to the input channel.
func (s *CliClient) Run(clientServiceCtx context.Context, input chan<- string) {

	reader := bufio.NewReader(os.Stdin)

	go func() {
		defer close(input)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("error reading the user input: %v", err)
				return
			}

			select {
			case <-clientServiceCtx.Done():
				return
			case input <- strings.TrimSpace(line):
			}
		}
	}()
}
