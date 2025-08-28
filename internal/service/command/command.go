package command

import (
	"context"
	"fmt"
	"log"
	"strings"

	syncservice "github.com/YenXXXW/clipboardSyncCliClient/internal/service/syncService"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
)

// CommandService is responsible for parsing user commands and delegating to other services.
type CommandService struct {
	clipService types.ClipService
	syncService *syncservice.SyncService
	input       chan string
}

// NewCommandService creates a new CommandService.
func NewCommandService(input chan string, syncService *syncservice.SyncService, clipService types.ClipService) *CommandService {
	return &CommandService{
		clipService: clipService,
		syncService: syncService,
		input:       input,
	}
}

// ProcessCommand parses the user input and executes the corresponding action.
func (s *CommandService) ProcessCommand(clientServiceCtx context.Context) {
	const (
		CmdCreate      = "/create"
		CmdLeave       = "/leave"
		CmdJoin        = "/join"
		CmdEnableSync  = "/enableSync"
		CmdDisableSync = "/disableSync"
	)
	const (
		Reset  = "\033[0m"
		Cyan   = "\033[36m"
		Yellow = "\033[33m"
	)

	fmt.Println()
	fmt.Println("You can use the following commands:")

	fmt.Printf("- Create a room:      %s\"%s\"%s\n", Cyan, CmdCreate, Reset)
	fmt.Printf("- Join a room:        %s\"%s <room_id>\"%s\n", Cyan, CmdJoin, Reset)
	fmt.Printf("- Leave a room:       %s\"%s\"%s\n", Cyan, CmdLeave, Reset)
	fmt.Printf("- Enable sync:        %s\"%s\"%s\n", Cyan, CmdEnableSync, Reset)
	fmt.Printf("- Disable sync:       %s\"%s\"%s\n", Cyan, CmdDisableSync, Reset)

	fmt.Println()
	fmt.Println(Yellow + "Waiting for your command..." + Reset)
	for {

		select {
		case commands, ok := <-s.input:
			if !ok {
				log.Println("stopping the command processing...")
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

			case CmdJoin:
				if len(args) < 1 {
					log.Println("Usage: /join <room_id>")
					continue
				}
				s.syncService.SubAndSyncUpdate(args[0])

			case CmdEnableSync:
				s.clipService.ToggleSyncEnable(true)

			case CmdDisableSync:
				s.clipService.ToggleSyncEnable(false)

			default:
				// If it's not a command, treat it as a clipboard update.
				fmt.Println("Please enter the correct command")
				fmt.Printf("to Create a room => \"%s\"\n", CmdCreate)
				fmt.Printf("to Join a room => \"%s\" <room_id>\n", CmdJoin)
				fmt.Printf("to Leave a room => \"%s\"\n", CmdLeave)
				fmt.Printf("to Enable Sync => \"%s\"\n", CmdEnableSync)
				fmt.Printf("to Disable Sync => \"%s\"\n", CmdDisableSync)
			}

		case <-clientServiceCtx.Done():
			fmt.Println("Process command stopped")
			return

		}
	}

}
