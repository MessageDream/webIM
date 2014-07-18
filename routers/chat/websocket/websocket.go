package websocket

import (
	"net/http"

	"github.com/MessageDream/webIM/models"
	"github.com/MessageDream/webIM/modules/base"
	"github.com/MessageDream/webIM/modules/log"
	"github.com/MessageDream/webIM/modules/middleware"
	"github.com/MessageDream/webIM/routers/chat"

	"github.com/gorilla/websocket"
)

const (
	WEBSOCKET base.TplName = "chat/websocket"
)

// Get method handles GET requests for WebSocketController.
func Get(ctx *middleware.Context) {
	// Safe check.
	uname := ctx.Query("uname")
	roomid := ctx.Query("roomid")
	if len(uname) == 0 {
		ctx.Redirect("/", 302)
		return
	}

	ctx.Data["IsWebSocket"] = true
	ctx.Data["UserName"] = uname
	ctx.Data["RoomID"] = roomid
	ctx.HTML(200, WEBSOCKET)
}

// Join method handles WebSocket requests for WebSocketController.
func Join(ctx *middleware.Context) {
	uname := ctx.Query("uname")
	roomid := ctx.Query("roomid")
	if len(uname) == 0 || len(roomid) == 0 {
		ctx.Redirect("/", 302)
		return
	}

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(ctx.Res, ctx.Req, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Error("Cannot setup WebSocket connection:", err)
		return
	}

	var room *chat.ChatRoom
	if chat.CheckRoom(roomid) {
		room = chat.GetRoom(roomid)
	} else {
		//Check roomid and then create room
		//check
		room = chat.NewChatRoom(roomid)
	}
	// Join chat room.
	room.Join(uname, ws)
	defer room.Leave(uname)

	// Message receive loop.
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			log.Error("%v", err)
			return
		}
		room.EnPublish(chat.NewEvent(models.EVENT_MESSAGE, models.User{Name: uname}, string(p)))
	}
}
