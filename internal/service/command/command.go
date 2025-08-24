package command

import (
	"context"
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

	go func() {
		for {

			select {
			case commands, ok := <-s.input:
				if !ok {
					log.Println("Error reading command from user in cli")
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
						return
					}
					s.syncService.SubAndSyncUpdate(args[0])

				case CmdEnableSync:
					s.clipService.ToggleSyncEnable(true)

				case CmdDisableSync:
					s.clipService.ToggleSyncEnable(false)

				default:
					// If it's not a command, treat it as a clipboard update.
					log.Println("Please enter the correct command")
					log.Printf("to Create a room => \"%s\"", CmdCreate)
					log.Printf("to Join a room => \"%s\" <room_id>", CmdJoin)
					log.Printf("to Leave a room => \"%s\"", CmdLeave)
					log.Printf("to Enable Sync => \"%s\"", CmdEnableSync)
					log.Printf("to Disable Sync => \"%s\"", CmdDisableSync)
				}

			case <-clientServiceCtx.Done():
				log.Println("Process command stopped")

			}
		}
	}()

}
