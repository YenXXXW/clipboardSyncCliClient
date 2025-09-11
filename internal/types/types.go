package types

import (
	"context"
)

/*  --- Services ---- */
type ClipService interface {
	RecieveUpdatesFromClipboardClient(context.Context)
	ProcessUpdates(*ClipboardUpdate)
	SendUpdate(context.Context, string)
	ToggleSyncEnable(bool)
}

// Notifier Service
type Notifier interface {
	Print(msg string)
}

type SyncService interface {
	SendUpdate(context.Context)
	LeaveRoom()
	CreateRoom()
	SubAndSyncUpdate(string) error
	GetRoomId() string
}

type Formatter interface {
	Info(string) string
	Success(string) string
	Error(string) string
	Warn(string) string
}

/* ---- Clients ------*/
// interface by the clip service for the clip client to be implemented
type ClipClient interface {
	ApplyUpdates(string)
	//fucntion that will give the newly update clipboard data
	NotifyUpdates(context.Context)
}

type CliClient interface {
	Run(ctx context.Context, input chan<- string)
}

// intface defince by sync service for sync client to be implemented
type SyncClient interface {
	SendUpdate(context.Context, string, string) error
	ReceiveUpdateAndSync(context.Context, string, string, chan *UpdateEvent) error
	LeaveRoom(context.Context, string, string) error
	CreateRoom(context.Context, string) (string, error)
}

/* -------------------  */
// internal type or domain type so that service does not depend on external dependencies's type
type ClipboardContent struct {
	Text string
}

type ClipboardUpdate struct {
	DeviceId string
	Content  *ClipboardContent
}

type Validate struct {
	Success bool
	Message string
}

type ValidateJoin struct {
	ValidateRoom Validate
	CheckClient  Validate
}

type UpdateEvent struct {
	ClipboardUpdate *ClipboardUpdate
	ValidateJoin    *ValidateJoin
}
