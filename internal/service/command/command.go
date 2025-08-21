package command

import (
	"context"
	"log"
	"strings"

	syncservice "github.com/YenXXXW/clipboardSyncCliClient/internal/service/syncService"
)

// CommandService is responsible for parsing user commands and delegating to other services.
type CommandService struct {
	syncService *syncservice.SyncService
	input       chan string
}

// NewCommandService creates a new CommandService.
func NewCommandService(input chan string, syncService *syncservice.SyncService) *CommandService {
	return &CommandService{
		syncService: syncService,
		input:       input,
	}
}

// ProcessCommand parses the user input and executes the corresponding action.
func (s *CommandService) ProcessCommand(ctx context.Context) {

	go func() {
		for commands := range s.input {

			parts := strings.Fields(commands)
			if len(parts) == 0 {
				continue
			}

			command := parts[0]
			args := parts[1:]

			switch command {
			case "/create":
				s.syncService.CreateRoom(ctx)
			case "/leave":
				s.syncService.LeaveRoom(ctx)
			case "/join":
				if len(args) < 1 {
					log.Println("Usage: /join <room_id>")
					return
				}
				s.syncService.SubAndSyncUpdate(ctx, args[0])
			default:
				// If it's not a command, treat it as a clipboard update.
				log.Println("Please enter the correct command")
			}
		}
	}()

}
