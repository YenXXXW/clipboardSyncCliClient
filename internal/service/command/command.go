package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
)

// CommandService is responsible for parsing user commands and delegating to other services.
type CommandService struct {
	notifier    types.Notifier
	formatter   types.Formatter
	clipService types.ClipService
	syncService types.SyncService
	input       chan string
}

// NewCommandService creates a new CommandService.
func NewCommandService(input chan string, syncService types.SyncService, clipService types.ClipService, formatter types.Formatter, notifier types.Notifier) *CommandService {
	return &CommandService{
		formatter:   formatter,
		notifier:    notifier,
		clipService: clipService,
		syncService: syncService,
		input:       input,
	}
}

const (
	CmdCreate      = "/create"
	CmdLeave       = "/leave"
	CmdJoin        = "/join"
	CmdEnableSync  = "/enableSync"
	CmdDisableSync = "/disableSync"
)

func (s *CommandService) PrintAvailableCommand() {
	s.notifier.Print("- Create a room:      " + s.formatter.Info(CmdCreate))
	s.notifier.Print("- Join a room:        " + s.formatter.Info(CmdJoin+" <room_id>"))
	s.notifier.Print("- Leave a room:       " + s.formatter.Info(CmdLeave))
	s.notifier.Print("- Enable sync:        " + s.formatter.Info(CmdEnableSync))
	s.notifier.Print("- Disable sync:       " + s.formatter.Info(CmdDisableSync))
}

// ProcessCommand parses the user input and executes the corresponding action.
func (s *CommandService) ProcessCommand(clientServiceCtx context.Context) {
	const (
		Reset  = "\033[0m"
		Cyan   = "\033[36m"
		Yellow = "\033[33m"
	)

	fmt.Println()
	s.notifier.Print(s.formatter.Warn("You can use the following commands:"))
	s.PrintAvailableCommand()

	fmt.Println()
	fmt.Println(Yellow + "Waiting for your command..." + Reset)
	for {

		select {
		case commands, ok := <-s.input:
			if !ok {
				s.notifier.Print(s.formatter.Info("Stopping the command processing"))
				return
			}

			parts := strings.Fields(commands)
			if len(parts) == 0 {
				continue
			}

			command := parts[0]
			args := parts[1:]

			switch command {
			case CmdCreate:
				s.syncService.CreateRoom()

			case CmdLeave:
				s.syncService.LeaveRoom()
				s.clipService.ToggleSyncEnable(false)

			case CmdJoin:
				if len(args) < 1 {
					s.notifier.Print(s.formatter.Warn("Usage: /join <room_id>"))
					continue
				}
				s.syncService.SubAndSyncUpdate(args[0])

			case CmdEnableSync:
				if roomID := s.syncService.GetRoomId(); roomID == "" {
					s.notifier.Print(s.formatter.Warn("Sync can be enabled only after joining a room"))
					continue
				}

				s.clipService.ToggleSyncEnable(true)

			case CmdDisableSync:
				if roomID := s.syncService.GetRoomId(); roomID == "" {
					s.notifier.Print(s.formatter.Warn("Please Join a room first"))
					continue
				}
				s.clipService.ToggleSyncEnable(false)

			default:
				// If it's not a command, treat it as a clipboard update.
				s.notifier.Print(s.formatter.Warn("Please enter the correct command"))
				s.PrintAvailableCommand()
			}

		case <-clientServiceCtx.Done():
			fmt.Println("Process command stopped")
			return

		}
	}

}
