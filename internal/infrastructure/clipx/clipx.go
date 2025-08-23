package clipxClient

import (
	"context"
	"fmt"
	"sync"

	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
	"golang.design/x/clipboard"
)

type clipxinfra struct {
	syncClient       types.SyncClient
	clipboardService types.ClipService
	isSyncingInbound bool //to prevent sedning of the clipbaord data to the server when we apply the remote change
	mutex            sync.Mutex
	deviceId         string
}

func NewClipxInfra(syncClient types.SyncClient, clipboardService types.ClipService, deviceId string) *clipxinfra {
	return &clipxinfra{
		syncClient:       syncClient,
		clipboardService: clipboardService,
		isSyncingInbound: false,
		deviceId:         deviceId,
	}
}

func (c *clipxinfra) Run() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	fmt.Println("the clipboard in ready")

}

func (c *clipxinfra) NotifyUpdates(ctx context.Context) {

	ch := clipboard.Watch(ctx, clipboard.FmtText)
	fmt.Println("watching the changes in clipboard")
	go func() {
		for data := range ch {
			c.clipboardService.Watch(string(data))
		}
	}()

}

func (c *clipxinfra) ApplyUpdates(content string) {

	clipboard.Write(clipboard.FmtText, []byte(content))
}
