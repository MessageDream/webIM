package socket

import (
	"encoding/json"
	"net"

	"github.com/MessageDream/webIM/models"
	"github.com/MessageDream/webIM/modules/log"
	"github.com/MessageDream/webIM/routers/chat"
)

type Client struct {
	UName string
	Room  *chat.ChatRoom
}

var clientMap map[net.Conn]*Client = make(map[net.Conn]*Client, 10)

func OnMessage(msg string, conn net.Conn) {
	client := clientMap[conn]
	if client.Room != nil {
		log.Info("%s", msg)
		client.Room.EnPublish(chat.NewEvent(models.EVENT_MESSAGE, models.User{Name: client.UName}, string(msg), client.Room.ID))
	} else {
		form := make(map[string]string, 2)
		json.Unmarshal([]byte(msg), &form)
		if len(form["uname"]) != 0 && len(form["roomid"]) != 0 {
			join(form["uname"], form["roomid"], conn)
		}
	}
}

func join(uname, roomid string, conn net.Conn) {
	var room *chat.ChatRoom
	if chat.CheckRoom(roomid) {
		room = chat.GetRoom(roomid)
	} else {
		//Check roomid and then create room
		//check
		room = chat.NewChatRoom(roomid)
	}
	client := clientMap[conn]
	client.Room = room
	client.UName = uname
	// Join chat room.
	room.Join(uname, nil, conn)
}

func OnConnected(conn net.Conn) {
	clientMap[conn] = &Client{}
}

func OnDisconnected(conn net.Conn) {
	clientMap[conn].Room.Leave(clientMap[conn].UName)
	delete(clientMap, conn)
}
