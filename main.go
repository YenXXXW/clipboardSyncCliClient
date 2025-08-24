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
	clipboardService "github.com/YenXXXW/clipboardSyncCliClient/internal/service/clipboard"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/service/command"
	syncservice "github.com/YenXXXW/clipboardSyncCliClient/internal/service/syncService"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
	"github.com/google/uuid"
)

func main() {
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

	userCliInputChan := make(chan string, 100)
	grpcCleint := infrastructure.NewGrpcClient("localhost:9000")
	cliClient := cliCleint.NewCliClient()

	log.Println("Reading data from the command line...")

	//update chan to send data between syncService and clipService
	updatesFromServerChan := make(chan *types.ClipboardUpdate, 100)
	localUpdatesChan := make(chan string, 100)

	//channel to send the data between the clipx infra and clipService
	localClipUpdatesChan := make(chan string, 100)

	clipxClient := clipxClient.NewClipxInfra(deviceId, localClipUpdatesChan)

	clipService := clipboardService.NewClipSyncService(deviceId, clipxClient, localUpdatesChan, localClipUpdatesChan)
	go clipService.RecieveUpdatesFromClipboardClient(clientServiceCtx)
	syncService := syncservice.NewSyncService(deviceId, "", clipService, grpcCleint, updatesFromServerChan, localUpdatesChan)

	commandService := command.NewCommandService(userCliInputChan, syncService, clipService)

	go syncService.SendUpdate(clientServiceCtx)
	clipxClient.Run()
	clipxClient.NotifyUpdates(clientServiceCtx)

	cliClient.Run(clientServiceCtx, userCliInputChan)
	commandService.ProcessCommand(clientServiceCtx)

	<-clientServiceCtx.Done()

}
