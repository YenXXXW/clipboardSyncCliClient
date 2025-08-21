package clipboardService

import (
	"context"
	"log"
	"sync"

	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
)

type ClipSyncService struct {
	clipClient       types.ClipClient
	clipSyncService  types.ClipSyncService
	isSyncingInbound bool
	deviceId         string
	mutex            sync.Mutex
}

func NewClipSyncService(clipSyncService types.ClipSyncService, deviceId string) *ClipSyncService {

	return &ClipSyncService{
		clipSyncService:  clipSyncService,
		isSyncingInbound: false,
		deviceId:         deviceId,
	}
}

func (c *ClipSyncService) Watch(data string) {

	//check if the change is initiated by the user or the program
	if !c.isSyncingInbound {
		if err := c.clipSyncService.SendUpdate(context.Background(), data); err != nil {
			log.Printf("failed to send update: %v", err)
		}
	}
}

// Identify the changes coming from remote and apply
func (c *ClipSyncService) ProcessUpdates(update *types.ClipboardUpdate) {
	//process and sync the update only if update not made by the same device
	if update.DeviceId != c.deviceId {
		c.clipClient.ApplyUpdates(update.Content.Text)
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
	return c.clipSyncService.SendUpdate(ctx, content)
}
