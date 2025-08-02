package clipboard

import (
// pb "github.com/YenXXXW/clipboardSyncCliClient/genproto/clipboardSync"
)

type ClipboardClient interface {
	SendUpdate(string) error
	SubScribeUpdate(string, string) error
}
