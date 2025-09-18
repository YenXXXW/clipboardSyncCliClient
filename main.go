package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	cliCleint "github.com/YenXXXW/clipboardSyncCliClient/internal/infrastructure/cli"
	clipxClient "github.com/YenXXXW/clipboardSyncCliClient/internal/infrastructure/clipx"
	infrastructure "github.com/YenXXXW/clipboardSyncCliClient/internal/infrastructure/grpcClient"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/infrastructure/notifier"
	clipboardService "github.com/YenXXXW/clipboardSyncCliClient/internal/service/clipboard"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/service/command"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/service/formatter"
	syncservice "github.com/YenXXXW/clipboardSyncCliClient/internal/service/syncService"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	clientServiceCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//listen for OS signals to gracefully shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	deviceId := uuid.New().String()

	terminalNotifier := notifier.NewTerminalNotifiter()
	formatter := formatter.NewFormatter()
	userCliInputChan := make(chan string, 100)
	grpcClient := infrastructure.NewGrpcClient(os.Getenv("SERVER_ADDR"), terminalNotifier, formatter)
	cliClient := cliCleint.NewCliClient()

	//update chan to send data between syncService and clipService
	updatesFromServerChan := make(chan *types.UpdateEvent, 100)
	localUpdatesChan := make(chan string, 100)

	defer close(updatesFromServerChan)
	defer close(localUpdatesChan)

	//channel to send the data between the clipx infra and clipService
	localClipUpdatesChan := make(chan string, 100)

	clipxClient := clipxClient.NewClipxInfra(deviceId, terminalNotifier, formatter, localClipUpdatesChan)

	clipService := clipboardService.NewClipSyncService(deviceId, clipxClient, localUpdatesChan, localClipUpdatesChan)
	go clipService.RecieveUpdatesFromClipboardClient(clientServiceCtx)
	syncService := syncservice.NewSyncService(formatter, terminalNotifier, deviceId, "", clipService, grpcClient, updatesFromServerChan, localUpdatesChan)

	commandService := command.NewCommandService(userCliInputChan, syncService, clipService, formatter, terminalNotifier)

	go syncService.SendUpdate(clientServiceCtx)
	clipxClient.Run()
	clipxClient.NotifyUpdates(clientServiceCtx)

	go commandService.ProcessCommand(clientServiceCtx)
	cliClient.Run(clientServiceCtx, userCliInputChan)

}
