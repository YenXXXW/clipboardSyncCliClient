package clipboard

import (
	"context"
	"fmt"

	"golang.design/x/clipboard"
)

type ClipSyncService struct {
	clipClient       ClipboardClient
	isSyncingInbound bool
}

func Init() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	fmt.Println("the program in working now")
}

func NewClipSyncService(clipClient ClipboardClient) *ClipSyncService {
	Init()
	return &ClipSyncService{
		clipClient:       clipClient,
		isSyncingInbound: false,
	}
}

func (c *ClipSyncService) Watch() {
	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for data := range ch {

		//check if the change is initiated by the user or the program
		if !c.isSyncingInbound {
			c.SendUpdate(string(data))
		}
	}
}

func (c *ClipSyncService) SendUpdate(content string) error {
	c.isSyncingInbound = true
	defer func() { c.isSyncingInbound = false }()
	return c.clipClient.SendUpdate(content)
}

func (c *ClipSyncService) SubAndSyncUpdate(deviceId, roomId string) error {
	return c.clipClient.SubScribeUpdate(deviceId, roomId)

}
