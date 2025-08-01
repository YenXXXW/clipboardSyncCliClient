package clipboard

import (
	"context"
	"fmt"

	"golang.design/x/clipboard"
)

type NewClipSyncService struct {
}

func Init() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	fmt.Println("the program in working now")
}

func Watch() {
	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for data := range ch {
		// print out clipboard data whenever it is changed
		println(string(data))
	}
}
