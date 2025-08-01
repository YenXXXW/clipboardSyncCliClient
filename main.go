package main

import (
	"github.com/YenXXXW/clipboardSyncCliClient/internal/grpcClient"
)

func main() {
	grpcClient.NewGrpcClient(":9000")
}
