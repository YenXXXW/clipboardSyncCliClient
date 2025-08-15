package infrastructure

import (
	"context"
	"fmt"
	"sync"

	clipboardService "github.com/YenXXXW/clipboardSyncCliClient/internal/service/clipboard"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
	"golang.design/x/clipboard"
)

type clipxinfra struct {
	syncClient       types.SyncClient
	clipboardService clipboardService.ClipSyncService
	isSyncingInbound bool //to prevent sedning of the clipbaord data to the server when we apply the remote change
	mutex            sync.Mutex
	deviceId         string
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
	for data := range ch {

		c.clipboardService.Watch(string(data))
	}

}

func (c *clipxinfra) ApplyUpdates(content string) {

	clipboard.Write(clipboard.FmtText, []byte(content))
}
