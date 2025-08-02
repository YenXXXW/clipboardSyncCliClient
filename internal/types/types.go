package types

type ClipboardService interface {
	Watch()
	SendUpdate()
	SubAndSyncUpdate()
}
