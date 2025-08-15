package syncservice

import (
	"context"
	"log"

	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"

	pb "github.com/YenXXXW/clipboardSyncCliClient/genproto/clipboardSync"
)

type SyncService struct {
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

func (s *SyncService) ReceiveUpdateAndSync(ctx context.Context, deviceId, roomId string, updateChan chan *pb.ClipboardUpdate) error {
	return s.syncClient.ReceiveUpdateAndSync(ctx, deviceId, roomId, updateChan)
}

func (c *SyncService) SubAndSyncUpdate(ctx context.Context, roomId string) error {
	streamCtx, cancel := context.WithCancel(ctx)
	c.cancelStream = cancel

	return c.syncClient.ReceiveUpdateAndSync(streamCtx, c.deviceId, roomId, c.incomingUpdatesFromServer)
}
