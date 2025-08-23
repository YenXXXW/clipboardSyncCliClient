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
	ctx, cancel := context.WithCancel(context.Background())
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

	updatesFromServerChan := make(chan *types.ClipboardUpdate, 100)
	syncService := syncservice.NewSyncService(deviceId, "", grpcCleint, updatesFromServerChan)

	clipService := clipboardService.NewClipSyncService(syncService, deviceId)

	clipxClient := clipxClient.NewClipxInfra(grpcCleint, clipService, deviceId)

	commandService := command.NewCommandService(userCliInputChan, syncService, clipService)
	clipxClient.Run()
	clipxClient.NotifyUpdates(ctx)
	cliClient.Run(ctx, userCliInputChan)
	commandService.ProcessCommand(ctx)

	<-ctx.Done()

}
