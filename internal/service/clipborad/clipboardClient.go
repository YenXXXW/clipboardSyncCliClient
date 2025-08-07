package clipboard

import (
	"context"

	pb "github.com/YenXXXW/clipboardSyncCliClient/genproto/clipboardSync"
)

type ClipboardClient interface {
	SendUpdate(context.Context, string, string) error
	ReceiveUpdateAndSync(context.Context, string, string, chan *pb.ClipboardUpdate) error
	LeaveRoom(context.Context, string, string) error
	CreateRoom(context.Context, string) (string, error)
}
