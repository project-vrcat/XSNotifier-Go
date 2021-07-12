package main

import (
	"fmt"
	"time"

	xsn "github.com/project-vrcat/XSNotifier-Go"
)

func main() {
	n, err := xsn.NewNotifier()
	if err != nil {
		panic(err)
	}
	defer n.Close()

	fmt.Println("Press Ctrl-C to exit.")

	for {
		n.Send(xsn.Message{
			Timeout:   3,
			Title:     "Example Notification!",
			Content:   time.Now().String(),
			SourceApp: "XSOverlay_Example_UDP",
			AudioPath: xsn.Audio_Default,
		})
		time.Sleep(time.Second * 10)
	}
}
