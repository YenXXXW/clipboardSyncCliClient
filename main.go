package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/YenXXXW/clipboardSyncCliClient/internal/grpcClient"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/service/cli"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/service/clipboard"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/service/command"
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

	grcClient := grpcClient.NewGrpcClient(":9000")
	clipSyncService := clipboard.NewClipSyncService(grcClient, uuid.NewString())
	go clipSyncService.Watch(ctx)
	commandService := command.NewCommandService(clipSyncService)
	cliService := cli.NewClipService()

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
