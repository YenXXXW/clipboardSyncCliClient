package clipboard

import (
	"context"
	"fmt"
	"log"
	"sync"

	pb "github.com/YenXXXW/clipboardSyncCliClient/genproto/clipboardSync"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
	"golang.design/x/clipboard"
)

type ClipSyncService struct {
	clipClient                types.SyncClient
	isSyncingInbound          bool
	incomingUpdatesFromServer chan *pb.ClipboardUpdate
	deviceId                  string
	roomId                    string
	mutex                     sync.Mutex
	cancelStream              context.CancelFunc
}

func Init() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	fmt.Println("the program in working now")
}

func NewClipSyncService(clipClient types.SyncClient, deviceId string) *ClipSyncService {
	Init()

	updateChan := make(chan *pb.ClipboardUpdate, 100)
	return &ClipSyncService{
		clipClient:                clipClient,
		isSyncingInbound:          false,
		incomingUpdatesFromServer: updateChan,
		deviceId:                  deviceId,
	}
}

func (c *ClipSyncService) Watch(ctx context.Context) {
	ch := clipboard.Watch(ctx, clipboard.FmtText)
	for data := range ch {

		//check if the change is initiated by the user or the program
		if !c.isSyncingInbound {
			c.SendUpdate(context.Background(), string(data))
		}
	}
}

func (c *ClipSyncService) SendUpdate(ctx context.Context, content string) error {

	c.mutex.Lock()
	c.isSyncingInbound = true
	c.mutex.Unlock()

	defer func() {

		c.mutex.Lock()
		c.isSyncingInbound = false
		c.mutex.Unlock()
	}()
	return c.clipClient.SendUpdate(ctx, c.deviceId, content)
}

func (c *ClipSyncService) SubAndSyncUpdate(ctx context.Context, roomId string) error {
	streamCtx, cancel := context.WithCancel(ctx)
	c.cancelStream = cancel

	return c.clipClient.ReceiveUpdateAndSync(streamCtx, c.deviceId, roomId, c.incomingUpdatesFromServer)
}

func (c *ClipSyncService) ProcessUpdates() {
	for update := range c.incomingUpdatesFromServer {
		//process and sync the update only if update not made by the same device
		if update.GetDeviceId() != c.deviceId {
			clipboard.Write(clipboard.FmtText, []byte(update.GetContent().GetText()))
		}
	}
}

func (c *ClipSyncService) CreateRoom(ctx context.Context) {
	roomId, err := c.clipClient.CreateRoom(ctx, c.deviceId)
	if err != nil {
		log.Printf("Error creating room %v", err)
		return
	}

	log.Println("roomId", roomId)
	c.roomId = roomId

}

func (c *ClipSyncService) LeaveRoom(ctx context.Context) {
	if c.cancelStream != nil {
		c.cancelStream()
	}
	c.clipClient.LeaveRoom(ctx, c.deviceId, c.roomId)
}
