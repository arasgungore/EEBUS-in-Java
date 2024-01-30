package ship

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"time"

	"github.com/LMF-DHBW/go_eebus/resources"

	"golang.org/x/net/websocket"
)

type SMEInstance struct {
	role            string
	connectionState string
	Connection      *websocket.Conn
	closeHandler    CloseHandler
	Ski             string
}

type CmiMessage struct {
	MessageType  int `json:"MessageType"`
	MessageValue int `json:"MessageValue"`
}

type Message struct {
	MessageType  int `xml:"MessageType"`
	MessageValue DataValue
}

type DataValue struct {
	Header  HeaderType             `xml:"header"`
	Payload resources.DatagramType `xml:"payload"`
}

type HeaderType struct {
	ProtocollId string `xml:"protocollId"`
}

const CMI_TIMEOUT = 10

func (SME *SMEInstance) RecieveTimeout(seconds int) []byte {
	ch := make(chan bool)
	msgch := make(chan []byte)
	go SME.RecieveOnce(func(msg []byte) {
		ch <- true
		msgch <- msg
	})
	time.AfterFunc(time.Second*time.Duration(seconds), func() { ch <- false })
	if !<-ch {
		msgch <- make([]byte, 0)
	}
	return <-msgch
}

func (SME *SMEInstance) StartCMI() {
	if SME.role == "server" {
		SME.connectionState = "CMI_STATE_SERVER_WAIT"
		msg := SME.RecieveTimeout(CMI_TIMEOUT)
		if len(msg) == 0 {
			SME.Connection.Close()
			return
		}
		SME.connectionState = "CMI_STATE_SERVER_EVALUATE"
		Message := CmiMessage{}
		resources.CheckError(json.Unmarshal(msg, &Message))
		if Message.MessageType != 0 || Message.MessageValue != 0 {
			defer SME.Connection.Close()
		}
		bytes, err := json.Marshal(CmiMessage{
			MessageType:  0,
			MessageValue: 0,
		})
		resources.CheckError(err)
		websocket.Message.Send(SME.Connection, bytes)
	} else {
		SME.connectionState = "CMI_STATE_CLIENT_SEND"
		bytes, err := json.Marshal(CmiMessage{
			MessageType:  0,
			MessageValue: 0,
		})
		resources.CheckError(err)
		websocket.Message.Send(SME.Connection, bytes)
		SME.connectionState = "CMI_STATE_CLIENT_WAIT"
		msg := SME.RecieveTimeout(CMI_TIMEOUT)
		if len(msg) == 0 {
			SME.Connection.Close()
			return
		}
		SME.connectionState = "CMI_STATE_CLIENT_EVALUATE"
		Message := CmiMessage{}
		resources.CheckError(json.Unmarshal(msg, &Message))
		if Message.MessageType != 0 || Message.MessageValue != 0 {
			SME.Connection.Close()
		}
	}
}

/* Use XML format
- Notify SPINE node that connection is active
- Recieve and send methods should be available
*/
func (SME *SMEInstance) RecieveOnce(handleFunc handler) {
	var msg []byte

	err := websocket.Message.Receive(SME.Connection, &msg)
	if err != nil {
		if err == io.EOF {
			return
		}
		SME.Connection.Close()
		SME.closeHandler(SME)
		return
	}
	handleFunc(msg)
}

func (SME *SMEInstance) Recieve(handleFunc dataHandler) {
	var msg []byte
	for {
		err := websocket.Message.Receive(SME.Connection, &msg)
		if err != nil {
			if err == io.EOF {
				break
			}
			SME.Connection.Close()
			SME.closeHandler(SME)
			break
		}
		Message := Message{}
		resources.CheckError(xml.Unmarshal(msg, &Message))
		if Message.MessageType == 2 {
			handleFunc(Message.MessageValue.Payload)
		}
	}
}

/* Sends messages in json format */
func (SME *SMEInstance) Send(payload resources.DatagramType) {
	bytes, err := xml.Marshal(Message{
		MessageType: 2,
		MessageValue: DataValue{
			Header: HeaderType{
				ProtocollId: "ee1.0",
			},
			Payload: payload,
		},
	})
	resources.CheckError(err)
	websocket.Message.Send(SME.Connection, bytes)
}
