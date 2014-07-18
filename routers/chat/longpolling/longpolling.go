package longpolling

import (
	"strconv"

	"github.com/MessageDream/webIM/models"
	"github.com/MessageDream/webIM/modules/base"
	"github.com/MessageDream/webIM/modules/log"
	"github.com/MessageDream/webIM/modules/middleware"
	"github.com/MessageDream/webIM/routers/chat"
)

const (
	LONGPOLLING base.TplName = "chat/longpolling"
)

// Join method handles GET requests.
func Join(ctx *middleware.Context) {
	// Safe check.
	uname := ctx.Query("uname")
	roomid := ctx.Query("roomid")

	if len(uname) == 0 || len(roomid) == 0 {
		ctx.Redirect("/", 302)
		return
	}
	var room *chat.ChatRoom
	if chat.CheckRoom(roomid) {
		room = chat.GetRoom(roomid)
	} else {
		//检查roomid格式并生成room
		//check
		room = chat.NewChatRoom(roomid)
	}

	room.Join(uname, nil)
	// Join chat room.
	//	chat.Join(uname, nil)

	ctx.Data["IsLongPolling"] = true
	ctx.Data["UserName"] = uname
	ctx.Data["RoomID"] = roomid
	ctx.HTML(200, LONGPOLLING)
}

// Post method handles receive messages requests.
func Post(ctx *middleware.Context) {
	uname := ctx.Query("uname")
	roomid := ctx.Query("roomid")
	content := ctx.Query("content")
	if len(uname) == 0 || len(content) == 0 || len(roomid) == 0 {
		log.Warn("%s", "One or more params are nil")
		return
	}
	if chat.CheckRoom(roomid) {
		room := chat.GetRoom(roomid)
		room.EnPublish(chat.NewEvent(models.EVENT_MESSAGE, models.User{Name: uname}, content))
	} else {
		log.Warn("There is no room's id is %s", roomid)
	}
}

// Fetch method handles fetch archives requests.
func Fetch(ctx *middleware.Context) {
	lastReceived, err := strconv.Atoi(ctx.Query("lastReceived"))
	roomid := ctx.Query("roomid")
	if err != nil || len(roomid) == 0 {
		return
	}

	if chat.CheckRoom(roomid) {
		events := models.GetEvents(roomid, int(lastReceived))
		if len(events) > 0 {
			ctx.Render.JSON(200, events)
			return
		}

		// Wait for new message(s).
		ch := make(chan bool)
		room := chat.GetRoom(roomid)
		room.PushBack(ch)
		<-ch

		events = models.GetEvents(roomid, int(lastReceived))
		ctx.Render.JSON(200, events)
		return
	}
	ctx.Render.JSON(200, nil)

}
