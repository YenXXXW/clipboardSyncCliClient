package types

import (
	"context"

	pb "github.com/YenXXXW/clipboardSyncCliClient/genproto/clipboardSync"
)

// intface defince by sync service for sync client to be implemented
type SyncClient interface {
	SendUpdate(context.Context, string, string) error
	ReceiveUpdateAndSync(context.Context, string, string, chan *ClipboardUpdate) error
	LeaveRoom(context.Context, string, string) error
	CreateRoom(context.Context, string) (string, error)
}

type ClipService interface {
	Watch(string)
	ProcessUpdates(*pb.ClipboardUpdate)
	SendUpdate(context.Context, string) error
}

// interface by the clip service for the clip client to be implemented
type ClipClient interface {
	ApplyUpdates(string)
	//fucntion that will give the newly update clipboard data
	NotifyUpdates(string)
}

type ClipSyncService interface {
	SendUpdate(context.Context, string) error
}

type CliClient interface {
	Run(ctx context.Context, input chan<- string)
}

// internal type or domain type so that service does not depend on external dependencies's type
type ClipboardContent struct {
	Text string
}

type ClipboardUpdate struct {
	DeviceId string
	Content  *ClipboardContent
}
