package cli

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"
)

type CliService struct {
	input  chan string
	cancel context.CancelFunc
}

func NewClipService() *CliService {

	return &CliService{}
}

func (s *CliService) Run(ctx context.Context, input chan<- string) {

	userInput := make(chan string, 100)

	reader := bufio.NewReader(os.Stdin)

	go func() {
		defer close(userInput)
		for {

			select {

			case <-ctx.Done():
				return

			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					log.Printf("error reading the user input %v", err)
					return
				}

				userInput <- strings.TrimSpace(line)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("cli input channel closed")
				return
			case line, ok := <-userInput:
				if !ok {
					log.Println("cli input channel closed (reader failed or exited)")
					return
				}
				input <- line

			}
		}
	}()

}
