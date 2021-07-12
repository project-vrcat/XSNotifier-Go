package xsn

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type XSNotifier struct {
	addr   string
	closed bool
	conn   net.Conn
	sync.Mutex
}

type Message struct {
	MessageType   int     `json:"messageType"`   // 1 = Notification Popup, 2 = MediaPlayer Information, will be extended later on.
	Index         int     `json:"index"`         // Only used for Media Player, changes the icon on the wrist.
	Timeout       float32 `json:"timeout"`       // How long the notification will stay on screen for in seconds
	Height        float32 `json:"height"`        // Height notification will expand to if it has content other than a title. Default is 175
	Opacity       float32 `json:"opacity"`       // Opacity of the notification, to make it less intrusive. Setting to 0 will set to 1.
	Volume        float32 `json:"volume"`        // Notification sound volume.
	AudioPath     string  `json:"audioPath"`     // File path to .ogg audio file. Can be "default", "error", or "warning". Notification will be silent if left empty.
	Title         string  `json:"title"`         // Notification title, supports Rich Text Formatting
	Content       string  `json:"content"`       // Notification content, supports Rich Text Formatting, if left empty, notification will be small.
	UseBase64Icon bool    `json:"useBase64Icon"` // Set to true if using Base64 for the icon image
	Icon          string  `json:"icon"`          // Base64 Encoded image, or file path to image. Can also be "default", "error", or "warning"
	SourceApp     string  `json:"sourceApp"`     // Somewhere to put your app name for debugging purposes
}

const (
	MessageType_NotificationPopup      = 1
	MessageType_MediaPlayerInformation = 2

	Icon_Default = "default"
	Icon_Error   = "error"
	Icon_Warning = "warning"

	Audio_Default = "default"
	Audio_Error   = "error"
	Audio_Warning = "warning"

	defaultVolume  = 0.7
	defaultOpacity = 1.0
	defaultTimeout = 3.0
	defaultHeight  = 120.0
)

func NewNotifier(port ...int) (*XSNotifier, error) {
	_port := 42069
	if len(port) > 0 {
		_port = port[0]
	}
	addr := fmt.Sprintf("127.0.0.1:%d", _port)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	return &XSNotifier{
		addr: addr,
		conn: conn,
	}, nil
}

func (o *XSNotifier) Send(msg Message) {
	if msg.MessageType < 1 {
		msg.MessageType = MessageType_NotificationPopup
	}
	if msg.Icon == "" {
		msg.Icon = Icon_Default
	}
	if (msg.AudioPath != "" && msg.Volume <= 0) || msg.Volume > 1 {
		msg.Volume = defaultVolume
	}
	if msg.Opacity > 1 || msg.Opacity <= 0 {
		msg.Opacity = defaultOpacity
	}
	if msg.Timeout > 60 || msg.Timeout <= 0 {
		msg.Timeout = defaultTimeout
	}
	if msg.Height <= 1 {
		msg.Height = defaultHeight
	}
	o.Lock()
	defer o.Unlock()
	if o.closed {
		return
	}
	data, _ := json.Marshal(msg)
	o.conn.Write(data)
}

func (o *XSNotifier) Close() {
	o.Lock()
	o.conn.Close()
	o.closed = true
	o.Unlock()
}
