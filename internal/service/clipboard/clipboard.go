package clipboardService

import (
	"context"
	"log"
	"sync"

	pb "github.com/YenXXXW/clipboardSyncCliClient/genproto/clipboardSync"
	syncservice "github.com/YenXXXW/clipboardSyncCliClient/internal/service/syncService"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
)

type ClipSyncService struct {
	clipClient       types.ClipClient
	synService       syncservice.SyncService
	isSyncingInbound bool
	deviceId         string
	mutex            sync.Mutex
}

func NewClipSyncService(syncService syncservice.SyncService, deviceId string) *ClipSyncService {

	return &ClipSyncService{
		synService:       syncService,
		isSyncingInbound: false,
		deviceId:         deviceId,
	}
}

func (c *ClipSyncService) Watch(data string) {

	//check if the change is initiated by the user or the program
	if !c.isSyncingInbound {
		if err := c.synService.SendUpdate(context.Background(), data); err != nil {
			log.Printf("failed to send update: %v", err)
		}
	}
}

// Identify the changes coming from remote and apply
func (c *ClipSyncService) ProcessUpdates(update *pb.ClipboardUpdate) {
	//process and sync the update only if update not made by the same device
	if update.GetDeviceId() != c.deviceId {
		c.clipClient.ApplyUpdates(update.GetContent().GetText())
	}
}

// Send the clipdata from local to server
func (c *ClipSyncService) SendUpdate(ctx context.Context, content string) error {

	c.mutex.Lock()
	c.isSyncingInbound = true
	c.mutex.Unlock()

	defer func() {

		c.mutex.Lock()
		c.isSyncingInbound = false
		c.mutex.Unlock()
	}()
	return c.synService.SendUpdate(ctx, content)
}
