package infrastructure

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/YenXXXW/clipboardSyncCliClient/genproto/clipboardSync"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type clipboardGrpcClient struct {
	client           pb.ClipSyncServiceClient
	terminalNotifier types.Notifier
	formatter        types.Formatter
}

func NewGrpcClient(addr string, terminalNotifer types.Notifier, formatter types.Formatter) *clipboardGrpcClient {
	//create a connection client
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}

	c := pb.NewClipSyncServiceClient(conn)
	clipboardClient := &clipboardGrpcClient{
		client:           c,
		terminalNotifier: terminalNotifer,
		formatter:        formatter,
	}
	return clipboardClient
}

func (c *clipboardGrpcClient) SendUpdate(ctx context.Context, deviceId, content string) error {

	reqContent := &pb.ClipboardContent{
		Text: content,
	}

	req := &pb.ClipboardUpdate{
		Content:  reqContent,
		DeviceId: deviceId,
	}

	if _, err := c.client.SendClipboardUpdate(ctx, req); err != nil {
		return err
	}

	return nil

}

func (c *clipboardGrpcClient) ReceiveUpdateAndSync(ctx context.Context, deviceId, roomId string, updateChan chan *types.UpdateEvent) error {

	req := &pb.SubscribeRequest{
		DeviceId: deviceId,
		RoomId:   roomId,
	}

	stream, err := c.client.SubscribeClipboardContentUpdate(ctx, req)
	if err != nil {
		fmt.Println("error")
		return err
	}
	c.terminalNotifier.Print(c.formatter.Success(fmt.Sprintf("Successfully Joined and Subscribed to room: %s", roomId)))

	go func() {
		//defer close(updateChan)

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("End of stream")
				break
			}
			if err != nil {
				return
			}

			var UpdateEvent types.UpdateEvent

			switch ev := resp.Event.(type) {

			case *pb.UpdateEvent_ValidateJoin:

				ValidateJoin := &types.ValidateJoin{
					ValidateRoom: types.Validate{
						Success: ev.ValidateJoin.ValidateRoom.Success,
						Message: ev.ValidateJoin.ValidateRoom.Message,
					},
					CheckClient: types.Validate{
						Success: ev.ValidateJoin.CheckClient.Success,
						Message: ev.ValidateJoin.CheckClient.Message,
					},
				}

				UpdateEvent.ValidateJoin = ValidateJoin
			case *pb.UpdateEvent_ClipboardUpdate:

				clipboardDataUpdate := &types.ClipboardUpdate{
					DeviceId: ev.ClipboardUpdate.GetDeviceId(),
					Content: &types.ClipboardContent{
						Text: ev.ClipboardUpdate.GetContent().GetText(),
					},
				}

				UpdateEvent.ClipboardUpdate = clipboardDataUpdate

			}

			select {
			case updateChan <- &UpdateEvent:
			case <-ctx.Done():
				return

			}

		}
	}()

	return nil

}

func (c *clipboardGrpcClient) CreateRoom(ctx context.Context, deviceId string) (string, error) {
	req := &pb.CreateRoomRequest{
		DeviceId: deviceId,
	}

	res, err := c.client.CreateRoom(ctx, req)
	if err != nil {
		log.Printf("Error creating room %v", err)
		return "", err
	}

	return res.GetRoomId(), nil

}

func (c *clipboardGrpcClient) LeaveRoom(ctx context.Context, deviceId, roomId string) error {
	req := &pb.LeaveRoomRequest{
		DeviceId: deviceId,
		RoomId:   roomId,
	}

	if _, err := c.client.LeaveRoom(ctx, req); err != nil {
		log.Printf("Error leaving room %v", err)
		return err
	}

	return nil

}
