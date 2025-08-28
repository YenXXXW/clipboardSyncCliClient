package clipxClient

import (
	"context"
	"fmt"
	"sync"

	"golang.design/x/clipboard"
)

type clipxinfra struct {
	localClipUpdatesChan chan string
	isSyncingInbound     bool //to prevent sedning of the clipbaord data to the server when we apply the remote change
	mutex                sync.Mutex
	deviceId             string
}

func NewClipxInfra(deviceId string, localClipUpdatesChan chan string) *clipxinfra {
	return &clipxinfra{
		localClipUpdatesChan: localClipUpdatesChan,
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
	fmt.Println("Watching for changes in clipboard")
	go func() {
		for data := range ch {
			c.localClipUpdatesChan <- string(data)
		}
	}()

}

func (c *clipxinfra) ApplyUpdates(content string) {

	clipboard.Write(clipboard.FmtText, []byte(content))
}
