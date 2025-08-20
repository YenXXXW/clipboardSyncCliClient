package syncservice

import (
	"context"
	"log"

	clipboardService "github.com/YenXXXW/clipboardSyncCliClient/internal/service/clipboard"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"

	pb "github.com/YenXXXW/clipboardSyncCliClient/genproto/clipboardSync"
)

type SyncService struct {
	clipbaorService           clipboardService.ClipSyncService
	syncClient                types.SyncClient
	deviceId                  string
	roomId                    string
	cancelStream              context.CancelFunc
	incomingUpdatesFromServer chan *pb.ClipboardUpdate
}

func NewSyncService(deviceId, roomId string, syncClient types.SyncClient, incomingUpdatesFromServer chan *pb.ClipboardUpdate) *SyncService {
	return &SyncService{
		syncClient:                syncClient,
		deviceId:                  deviceId,
		roomId:                    roomId,
		incomingUpdatesFromServer: incomingUpdatesFromServer,
	}
}

func (s *SyncService) SendUpdate(ctx context.Context, content string) error {
	return s.syncClient.SendUpdate(ctx, s.deviceId, content)
}

func (s *SyncService) CreateRoom(ctx context.Context) {
	roomId, err := s.syncClient.CreateRoom(ctx, s.deviceId)
	if err != nil {
		log.Printf("Error creating room %v", err)
		return
	}

	log.Println("roomId", roomId)
	s.roomId = roomId

}

func (s *SyncService) LeaveRoom(ctx context.Context) {
	if s.cancelStream != nil {
		s.cancelStream()
	}
	s.syncClient.LeaveRoom(ctx, s.deviceId, s.roomId)
}

func (c *SyncService) SubAndSyncUpdate(ctx context.Context, roomId string) error {
	streamCtx, cancel := context.WithCancel(ctx)
	c.cancelStream = cancel

	go func() {
		for {
			select {
			case <-streamCtx.Done():
				log.Println("Stopping update processor: context canceled")
				return
			case update, ok := <-c.incomingUpdatesFromServer:
				if !ok {
					log.Println("Stopping update processor: channel closed")
					return
				}
				log.Printf("Received update from server: %s", update.GetContent().GetText())
				// Process the updates from the incomingudpates channel sent by grpc client and apply it to the clipboard
				c.clipbaorService.ProcessUpdates(update)
			}
		}
	}()

	if err := c.syncClient.ReceiveUpdateAndSync(streamCtx, c.deviceId, roomId, c.incomingUpdatesFromServer); err != nil {
		log.Printf("failed to subscribe to updates: %v", err)
		return err
	}

	return nil
}
