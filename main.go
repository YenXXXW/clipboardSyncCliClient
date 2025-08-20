package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	infrastructure "github.com/YenXXXW/clipboardSyncCliClient/internal/infrastructure/grpcClient"
	clipboardService "github.com/YenXXXW/clipboardSyncCliClient/internal/service/clipboard"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/service/command"
	syncservice "github.com/YenXXXW/clipboardSyncCliClient/internal/service/syncService"
	"github.com/google/uuid"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	commandInput := make(chan string, 100)
	deviceId := uuid.NewString()

	grcClient := infrastructure.NewGrpcClient(":9000")
	clipboardService := syncservice.NewSyncService()
	clipSyncService := clipboardService.NewClipSyncService(clipboardService, deviceId)
	go clipSyncService.Watch(ctx)
	commandService := command.NewCommandService(commandInput, clipSyncService)

	userInputChan := make(chan string, 100)
	go cliService.Run(ctx, userInputChan)

	log.Println("application started. enter /create to create a room, /join <room_id> to join a room")

	for {
		select {
		case <-ctx.Done():
			log.Println("application shutting down")
			return
		case input := <-userInputChan:
			commandService.ProcessCommand(ctx, input)
		}
	}
}
