package clipxClient

import (
	"context"
	"sync"

	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
	"golang.design/x/clipboard"
)

type clipxinfra struct {
	localClipUpdatesChan chan string
	notifier             types.Notifier
	formatter            types.Formatter
	isSyncingInbound     bool //to prevent sedning of the clipbaord data to the server when we apply the remote change
	mutex                sync.Mutex
	deviceId             string
}

func NewClipxInfra(deviceId string, notifier types.Notifier, formatter types.Formatter, localClipUpdatesChan chan string) *clipxinfra {
	return &clipxinfra{
		localClipUpdatesChan: localClipUpdatesChan,
		notifier:             notifier,
		formatter:            formatter,
		isSyncingInbound:     false,
		deviceId:             deviceId,
	}
}

func (c *clipxinfra) Run() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
}

func (c *clipxinfra) NotifyUpdates(clientServiceCtx context.Context) {

	ch := clipboard.Watch(clientServiceCtx, clipboard.FmtText)
	c.notifier.Print(c.formatter.Info("Watching for changes in clipboard"))
	go func() {
		defer c.notifier.Print(c.formatter.Error("Stopped Watching Clipboard..."))
		for data := range ch {
			c.localClipUpdatesChan <- string(data)
		}
	}()

}

func (c *clipxinfra) ApplyUpdates(content string) {

	clipboard.Write(clipboard.FmtText, []byte(content))
}
