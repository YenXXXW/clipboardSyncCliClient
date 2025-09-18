package cliCleint

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/chzyer/readline"
)

type CliClient struct {
}

func NewCliClient() *CliClient {

	return &CliClient{}
}

// Run starts reading user input from the command line continuously and sends it to the input channel.
func (s *CliClient) Run(clientServiceCtx context.Context, input chan<- string) {

	rl, err := readline.New("> ")
	if err != nil {
		fmt.Println("Unable to read the input from terminal")
		panic(err)
	}
	rlChan := make(chan string, 10)
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer close(rlChan)
		defer wg.Done()
		for {
			line, err := rl.Readline()
			if err != nil {
				log.Printf("error reading the user input: %v", err)
				return
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			rlChan <- line

		}
	}()

	wg.Add(1)
	go func() {
		defer close(input)
		defer wg.Done()

		for {

			select {
			case <-clientServiceCtx.Done():
				log.Println("clipboard client stopped")
				return
			case line, ok := <-rlChan:
				if !ok {
					return
				}
				input <- line
			}
		}
	}()

	wg.Wait()
	rl.Close()
}
