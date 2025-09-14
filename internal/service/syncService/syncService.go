package syncservice

import (
	"context"
	"fmt"
	"time"

	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
)

type SyncService struct {
	formatter                 types.Formatter
	localUpdatesChan          chan string
	clipboardService          types.ClipService
	syncClient                types.SyncClient
	deviceId                  string
	roomId                    string
	cancelStream              context.CancelFunc
	incomingUpdatesFromServer chan *types.UpdateEvent
	infoLogger                types.Notifier
}

func NewSyncService(formatter types.Formatter, infoLogger types.Notifier, deviceId, roomId string, clipboardService types.ClipService, syncClient types.SyncClient, incomingUpdatesFromServer chan *types.UpdateEvent, LocalUpdateChan chan string) *SyncService {
	return &SyncService{
		formatter:                 formatter,
		infoLogger:                infoLogger,
		localUpdatesChan:          LocalUpdateChan,
		clipboardService:          clipboardService,
		syncClient:                syncClient,
		deviceId:                  deviceId,
		roomId:                    roomId,
		incomingUpdatesFromServer: incomingUpdatesFromServer,
	}
}

// INFO: Fuction to receive the updates from the local updates channel
func (s *SyncService) SendUpdate(clientServiceCtx context.Context) {
	for {
		select {

		case content, ok := <-s.localUpdatesChan:
			if !ok {
				return
			}
			s.sendRpcWithTimeout(content)

		case <-clientServiceCtx.Done():
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
		s.infoLogger.Print(s.formatter.Error(fmt.Sprintf("Error creating room %v", err)))
		return
	}

	s.SubAndSyncUpdate(roomId)

	successMsg := s.formatter.Success("Room created successfully")
	roomIdMsg := s.formatter.Info(fmt.Sprintf("room id - %s", roomId))
	s.infoLogger.Print(successMsg)
	s.infoLogger.Print(roomIdMsg)
}

func (s *SyncService) LeaveRoom() {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if s.cancelStream != nil {
		s.cancelStream()
	}
	//disable the sync when the user leaves the room

	if s.roomId == "" {
		errorMsg := s.formatter.Error("You are not in a room")
		s.infoLogger.Print(errorMsg)
		return
	}
	s.syncClient.LeaveRoom(ctx, s.deviceId, s.roomId)
	s.roomId = ""
	s.infoLogger.Print(s.formatter.Success("Left Rooom Successfully"))
}

func (c *SyncService) SubAndSyncUpdate(roomId string) error {
	streamCtx, cancel := context.WithCancel(context.Background())
	c.cancelStream = cancel

	validatedAlready := false

	go func() {
		for {
			select {
			case <-streamCtx.Done():
				return
			case updateEvent, ok := <-c.incomingUpdatesFromServer:
				if !ok {
					return
				}

				if updateEvent.ValidateJoin == nil {
					// only perform the the below opreation once for the first response from the server
					update := updateEvent.ClipboardUpdate

					c.clipboardService.ProcessUpdates(update)
				} else {
					validateResult := updateEvent.ValidateJoin
					if !validateResult.ValidateRoom.Success {
						c.infoLogger.Print(c.formatter.Error(validateResult.ValidateRoom.Message))
						return
					} else if !validateResult.CheckClient.Success {
						c.infoLogger.Print(c.formatter.Error(validateResult.CheckClient.Message))
						return
					} else {
						if !validatedAlready {
							c.roomId = roomId
							c.clipboardService.ToggleSyncEnable(true)
							c.infoLogger.Print(c.formatter.Info("Successfully Joined the room"))
							validatedAlready = true
						}
					}
				}

				// Process the updates from the incomingudpates channel sent by grpc client and apply it to the clipboard
			}
		}
	}()

	if err := c.syncClient.ReceiveUpdateAndSync(streamCtx, c.deviceId, roomId, c.incomingUpdatesFromServer); err != nil {
		c.infoLogger.Print(c.formatter.Error(fmt.Sprintf("failed to subscribe to updates: %v", err)))
		cancel()
		c.roomId = ""
		c.clipboardService.ToggleSyncEnable(false)
		return err
	}

	return nil
}

func (c *SyncService) GetRoomId() string {
	return c.roomId
}
