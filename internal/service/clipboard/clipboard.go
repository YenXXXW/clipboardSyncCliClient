package clipboardService

import (
	"context"
	"log"
	"sync"

	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
)

type ClipSyncService struct {
	syncEnabled          bool
	clipClient           types.ClipClient
	isSyncingInbound     bool
	deviceId             string
	mutex                sync.Mutex
	localUpdatesChan     chan string
	localClipUpdatesChan chan string
}

func NewClipSyncService(deviceId string, clipClient types.ClipClient, LocalUpdateChan chan string, localClipUpdatesChan chan string) *ClipSyncService {

	return &ClipSyncService{
		clipClient:           clipClient,
		syncEnabled:          false,
		isSyncingInbound:     false,
		deviceId:             deviceId,
		localUpdatesChan:     LocalUpdateChan,
		localClipUpdatesChan: localClipUpdatesChan,
	}
}

func (c *ClipSyncService) RecieveUpdatesFromClipboardClient(clientServiceCtx context.Context) {

	for {
		select {
		case <-clientServiceCtx.Done():
			return

		case data := <-c.localClipUpdatesChan:

			//check if the change is initiated by the user or the program
			if !c.isSyncingInbound && c.syncEnabled {
				//use the local funciton to apply business rules
				c.SendUpdate(context.Background(), data)
			}

		}
	}
}

// Identify the changes coming from remote and apply
func (c *ClipSyncService) ProcessUpdates(update *types.ClipboardUpdate) {
	//process and sync the update only if update not made by the same device
	log.Println("insdie process Updates function")
	log.Println(update.DeviceId, c.deviceId, c.syncEnabled)
	if update.DeviceId != c.deviceId && c.syncEnabled {
		c.clipClient.ApplyUpdates(update.Content.Text)
	}
}

// Send the clipdata from local to server by applying business rules
func (c *ClipSyncService) SendUpdate(ctx context.Context, content string) {

	c.mutex.Lock()
	c.isSyncingInbound = true
	c.mutex.Unlock()

	defer func() {
		c.mutex.Lock()
		c.isSyncingInbound = false
		c.mutex.Unlock()
	}()

	c.localUpdatesChan <- content
}

func (c *ClipSyncService) ToggleSyncEnable(state bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	log.Println("clipboard sync state", state)
	c.syncEnabled = state
}
