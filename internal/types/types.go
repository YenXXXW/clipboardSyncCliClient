package types

import (
	"context"

	pb "github.com/YenXXXW/clipboardSyncCliClient/genproto/clipboardSync"
)

type SyncClient interface {
	SendUpdate(context.Context, string, string) error
	ReceiveUpdateAndSync(context.Context, string, string, chan *pb.ClipboardUpdate) error
	LeaveRoom(context.Context, string, string) error
	CreateRoom(context.Context, string) (string, error)
}

type ClipClient interface {
	ApplyUpdates(string)
	//fucntion that will give the newly update clipboard data
	NotifyUpdates(string)
}
