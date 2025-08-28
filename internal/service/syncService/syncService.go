package syncservice

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
)

type SyncService struct {
	localUpdatesChan          chan string
	clipboardService          types.ClipService
	syncClient                types.SyncClient
	deviceId                  string
	roomId                    string
	cancelStream              context.CancelFunc
	incomingUpdatesFromServer chan *types.ClipboardUpdate
	infoLogger                types.Notifier
}

func NewSyncService(infoLogger types.Notifier, deviceId, roomId string, clipboardService types.ClipService, syncClient types.SyncClient, incomingUpdatesFromServer chan *types.ClipboardUpdate, LocalUpdateChan chan string) *SyncService {
	return &SyncService{
		infoLogger:                infoLogger,
		localUpdatesChan:          LocalUpdateChan,
		clipboardService:          clipboardService,
		syncClient:                syncClient,
		deviceId:                  deviceId,
		roomId:                    roomId,
		incomingUpdatesFromServer: incomingUpdatesFromServer,
	}
}

func (s *SyncService) SendUpdate(clientServiceCtx context.Context) {
	for {
		select {

		case content, ok := <-s.localUpdatesChan:
			if !ok {
				log.Println("LocalUpdatesChan was closed. Stopping.")
				return
			}
			s.sendRpcWithTimeout(content)

		case <-clientServiceCtx.Done():
			fmt.Println("Sync service stopped")
			return
		}

	}
}

func (s *SyncService) sendRpcWithTimeout(content string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.syncClient.SendUpdate(ctx, s.deviceId, content)
}

func (s *SyncService) CreateRoom() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roomId, err := s.syncClient.CreateRoom(ctx, s.deviceId)

	if err != nil {
		log.Printf("Error creating room %v", err)
		return
	}

	s.SubAndSyncUpdate(roomId)

	s.infoLogger.Success("Room created successfully")
	s.infoLogger.Info(fmt.Sprintf("room id - %s", roomId))
}

func (s *SyncService) LeaveRoom() {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if s.cancelStream != nil {
		s.cancelStream()
	}
	//disable the sync when the user leaves the room

	if s.roomId == "" {
		s.infoLogger.Error("You are not in a room")
		return
	}
	s.syncClient.LeaveRoom(ctx, s.deviceId, s.roomId)
	s.roomId = ""
	fmt.Println("roomid", s.roomId)
	s.infoLogger.Success("Left Rooom Successfully")
}

func (c *SyncService) SubAndSyncUpdate(roomId string) error {
	streamCtx, cancel := context.WithCancel(context.Background())
	c.cancelStream = cancel

	go func() {
		for {
			select {
			case <-streamCtx.Done():
				return
			case update, ok := <-c.incomingUpdatesFromServer:
				if !ok {
					return
				}
				log.Printf("Received update from server: %v", update)
				// Process the updates from the incomingudpates channel sent by grpc client and apply it to the clipboard
				c.clipboardService.ProcessUpdates(update)
			}
		}
	}()

	c.clipboardService.ToggleSyncEnable(true)
	c.roomId = roomId

	if err := c.syncClient.ReceiveUpdateAndSync(streamCtx, c.deviceId, roomId, c.incomingUpdatesFromServer); err != nil {
		log.Printf("failed to subscribe to updates: %v", err)
		cancel()
		c.clipboardService.ToggleSyncEnable(false)
		return err
	}

	return nil
}
