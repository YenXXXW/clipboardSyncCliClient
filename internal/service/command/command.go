package command

import (
	"context"
	"log"
	"strings"

	"github.com/YenXXXW/clipboardSyncCliClient/internal/service/clipboard"
)

// CommandService is responsible for parsing user commands and delegating to other services.
type CommandService struct {
	clipSyncService *clipboard.ClipSyncService
}

// NewCommandService creates a new CommandService.
func NewCommandService(clipSyncService *clipboard.ClipSyncService) *CommandService {
	return &CommandService{
		clipSyncService: clipSyncService,
	}
}

// ProcessCommand parses the user input and executes the corresponding action.
func (s *CommandService) ProcessCommand(ctx context.Context, input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "/create":
		s.clipSyncService.CreateRoom(ctx)
	case "/leave":
		s.clipSyncService.LeaveRoom(ctx)
	case "/join":
		if len(args) < 1 {
			log.Println("Usage: /join <room_id>")
			return
		}
		s.clipSyncService.SubAndSyncUpdate(ctx, args[0])
	default:
		// If it's not a command, treat it as a clipboard update.
		log.Println("Please enter the correct command")
	}
}
